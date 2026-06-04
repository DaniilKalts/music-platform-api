#!/bin/bash
set -e

# Load environment variables from .env
if [ -f .env ]; then
    echo "Loading environment from .env..."
    while IFS='=' read -r key value; do
        [[ $key =~ ^#.* ]] && continue
        [[ -z $key ]] && continue
        value=$(echo "$value" | tr -d '\r')
        export "$key=$value"
    done < .env
fi

# Configuration with defaults
S3_ENDPOINT=${S3_ENDPOINT:-"localhost:9000"}
S3_ACCESS_KEY=${S3_ACCESS_KEY:-"admin"}
S3_SECRET_KEY=${S3_SECRET_KEY:-"password"}
S3_BUCKET=${S3_BUCKET:-"tracks"}

# Ensure endpoint has protocol
if [[ ! $S3_ENDPOINT == http* ]]; then
    S3_ENDPOINT="http://$S3_ENDPOINT"
fi

echo "--- S3 Seeding Started ---"
echo "Endpoint: $S3_ENDPOINT"
echo "Bucket:   $S3_BUCKET"

# ----------------------------------------------------------------------------
# Canonical, world-famous public-domain recordings (verified: HTTP 200, real
# MP3 audio). All sourced from the Internet Archive; the compositions are public
# domain (composers died 200+ years ago).
#
# The object keys here match the file_url values seeded by
# database/migrations/00011_seed_data.sql — keep the two files in sync.
#
# Format is "object-key.mp3|download-url". The /download/ path 302-redirects to
# a live Archive.org node; curl -L follows it. Spaces in file names are %20.
# ----------------------------------------------------------------------------
IA="https://archive.org/download"
FS="$IA/The_Four_Seasons_Vivaldi-10361/John_Harrison_with_the_Wichita_State_University_Chamber_Players"
TRACKS=(
  # Vivaldi — "The Four Seasons" (John Harrison / Wichita State University Chamber Players)
  "vivaldi-spring.mp3|${FS}_-_01_-_Spring_Mvt_1_Allegro.mp3"
  "vivaldi-summer.mp3|${FS}_-_06_-_Summer_Mvt_3_Presto.mp3"
  "vivaldi-autumn.mp3|${FS}_-_07_-_Autumn_Mvt_1_Allegro.mp3"
  "vivaldi-winter.mp3|${FS}_-_10_-_Winter_Mvt_1_Allegro_non_molto.mp3"
  # Beethoven
  "beethoven-moonlight.mp3|$IA/MoonlightSonata_845/Sonata_no_14_in_c_sharp_minor_moonlight_op_27_no_2_Iii.Presto.mp3"
  "beethoven-fur-elise.mp3|$IA/WoO59PocoMotoBagatelleInAMinorFurElise/FrEliseWoo59.mp3"
  # Mozart
  "mozart-eine-kleine-nachtmusik.mp3|$IA/SerenadeNo.13EineKleineNachtmusikK.525/Serenade%20No.%2013%20Eine%20Kleine%20Nachtmusik,%20K.525%20-%20I.%20Allegro.mp3"
  "mozart-turkish-march.mp3|$IA/MozartTurkishMarch/Mozart-TurkishMarch.mp3"
  # Bach
  "bach-toccata-fugue.mp3|$IA/ToccataAndFugueInDMinorBWV565_201310/Toccata%20and%20Fugue%20in%20D%20Minor,%20BWV%20565%20-%20I.%20Toccata.mp3"
  "bach-air-on-g-string.mp3|$IA/Bach-airOnTheGString/LaMusicaClasicaMasRelajanteDelMundo-Bach-AirOnTheGString.mp3"
  # Pachelbel
  "pachelbel-canon.mp3|$IA/PachelbelsCanoninD/Canon_in_D_Piano.mp3"
)

# Download every track to a temp dir on the host. Doing the fetch here (rather
# than streaming inside the mc container) gives us reliable TLS + redirects and
# lets us verify each file is real audio before uploading.
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${#TRACKS[@]} canonical tracks..."
for entry in "${TRACKS[@]}"; do
    key="${entry%%|*}"
    url="${entry#*|}"
    echo "  -> $key"
    curl -fSL --retry 3 --retry-delay 2 --max-time 300 -o "$TMPDIR/$key" "$url"

    # Sanity-check: non-empty and actually an MPEG audio stream.
    if [ ! -s "$TMPDIR/$key" ]; then
        echo "ERROR: $key downloaded empty from $url" >&2
        exit 1
    fi
    if command -v file &> /dev/null && ! file "$TMPDIR/$key" | grep -qiE "audio|mpeg|mp3"; then
        echo "ERROR: $key does not look like audio (got: $(file -b "$TMPDIR/$key"))" >&2
        exit 1
    fi
done

# Build the mc command sequence. Files live in /seed inside the container
# (or the same temp dir when mc runs locally).
build_commands() {
    local srcdir="$1"
    echo "set -e"
    echo "mc alias set myminio \"$S3_ENDPOINT\" \"$S3_ACCESS_KEY\" \"$S3_SECRET_KEY\""
    echo "echo \"Ensuring bucket '$S3_BUCKET' exists...\""
    echo "mc mb --ignore-existing myminio/\"$S3_BUCKET\""
    echo "echo \"Setting public download policy...\""
    echo "mc anonymous set download myminio/\"$S3_BUCKET\" || echo \"Warning: Could not set anonymous policy.\""
    for entry in "${TRACKS[@]}"; do
        local key="${entry%%|*}"
        echo "echo \"  uploading $key...\""
        echo "mc cp \"$srcdir/$key\" myminio/\"$S3_BUCKET/$key\""
    done
}

if command -v mc &> /dev/null; then
    # Run locally if mc exists
    bash -c "$(build_commands "$TMPDIR")"
else
    # Run inside a single docker container; mount the downloaded tracks at /seed
    echo "MinIO Client (mc) not found locally, using Docker..."
    docker run --rm -i --network host \
        -v "$TMPDIR:/seed:ro" \
        --entrypoint=/bin/sh minio/mc -c "$(build_commands "/seed")"
fi

echo "--- S3 Seeding Completed Successfully (${#TRACKS[@]} tracks) ---"
