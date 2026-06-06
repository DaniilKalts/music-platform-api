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

  # ── Ragtime — Scott Joplin ──
  "joplin-maple-leaf-rag.mp3|$IA/ScottJoplinRagtime/ScottJoplinMapleLeaf.mp3"
  "joplin-the-entertainer.mp3|$IA/joplin_ragtime_jop_20_scott_jop/joplin_ragtime_jop_01_the_enter.mp3"
  "joplin-pine-apple-rag.mp3|$IA/joplin_ragtime_jop_20_scott_jop/joplin_ragtime_jop_02_pine_appl.mp3"
  "joplin-the-ragtime-dance.mp3|$IA/joplin_ragtime_jop_20_scott_jop/joplin_ragtime_jop_04_the_ragti.mp3"
  "joplin-elite-syncopations.mp3|$IA/joplin_ragtime_jop_20_scott_jop/joplin_ragtime_jop_07_elite_syn.mp3"
  # ── Early Jazz ──
  "odjb-livery-stable-blues.mp3|$IA/OriginalDixielandJassBand/OriginalDixielandJassBand-LiveryStableBlues.mp3"
  "odjb-tiger-rag.mp3|$IA/78_9903-Tiger-Rag/9903-Tiger-Rag.mp3"
  "funny-jas-band-from-dixieland.mp3|$IA/fjasband1916/fjasband1916.mp3"
  # ── Early Blues ──
  "mamie-smith-crazy-blues.mp3|$IA/MamieSmithHerJazzHounds/MamieSmithHerJazzHounds-CrazyBlues.mp3"
  "handy-st-louis-blues.mp3|$IA/OriginalDixielandJazzBandwithAlBernard/OriginalDixielandJazzBandwithAlBernard-StLouisBlues.mp3"
  "handy-memphis-blues.mp3|$IA/Victor_Military_Band-The_Memphis_Blues/Victor_Military_Band-The_Memphis_Blues/Victor_Military_Band-The_Memphis_Blues.mp3"
  # ── Marches & Opera ──
  "sousa-stars-and-stripes-forever.mp3|$IA/JOHNPHILIPSOUSAMarches-NEWTRANSFER/01.StarsAndStripesForever.mp3"
  "sousa-washington-post-march.mp3|$IA/JOHNPHILIPSOUSAMarches-NEWTRANSFER/02.WashingtonPostMarch.mp3"
  "caruso-o-sole-mio.mp3|$IA/Caruso_part1/Caruso-OSoleMio.mp3"
  # ── Folk / Americana ──
  "swing-low-sweet-chariot.mp3|$IA/SwingLowSweetChariot_201609/Swing%20Low%20Sweet%20Chariot.mp3"
  "old-folks-at-home.mp3|$IA/us-oldfolks/us-oldfolks.mp3"
  # ── More iconic Classical ──
  "beethoven-symphony-5.mp3|$IA/SymphonyNo.5Opus67/Symphony%20No.%205%20-%20Opus%2067%2C%201st%20Movement.mp3"
  "beethoven-symphony-9-ode-to-joy.mp3|$IA/beethoven9/beethoven-9-04-concertgebouw-klemperer-1956-16048.mp3"
  "bach-brandenburg-3.mp3|$IA/BrandenburgConcertoNo.3InGMajor/BrandenburgConcertoNo.3InGMajor-I.Allegro.mp3"
  "bach-cello-suite-1-prelude.mp3|$IA/15SuiteNo.4EnMiBemolMajeur/01%20Suite%20no.%201%2C%20en%20sol%20majeur%2C%20pour.mp3"
  "chopin-nocturne-9-2.mp3|$IA/Chopin-NocturneOp.9No.2/20120420_Chopin_Nocturne_op9-2_amplified.mp3"
  "grieg-mountain-king.mp3|$IA/16EdvardGriegPeerGyntInTheHallOfTheMountainKing1875/16%20Edvard%20Grieg%20-%20Peer%20Gynt%20In%20The%20Hall%20Of%20The%20Mountain%20King%2C1875.mp3"
  "saint-saens-danse-macabre.mp3|$IA/DanseMacabreOp.40/Danse%20Macabre%2C%20Op.%2040.mp3"
  "mozart-symphony-40.mp3|$IA/SymphonyNo.40InGMinor/SymphonyNo.40InGMinorKv.550-I.AllegroModerato.mp3"
  "rimsky-flight-bumblebee.mp3|$IA/FlightOfTheBumblebee_201310/Flight%20of%20the%20Bumblebee.mp3"
  "handel-hallelujah.mp3|$IA/HallelujahChorusMessiah/EDIS-SRP-0195-06_hallelujah_chorus.mp3"
  "debussy-clair-de-lune.mp3|$IA/ClairDeLune_182/Debussy_Clair_de_Lune.mp3"
  "schubert-ave-maria.mp3|$IA/FranzSchubertAveMaria/Franz%20Schubert_%20Ave%20Maria.mp3"
)

# Download every track to a temp dir on the host. Doing the fetch here (rather
# than streaming inside the mc container) gives us reliable TLS + redirects and
# lets us verify each file is real audio before uploading.
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Fetch one track into $TMPDIR/$key, verified as real audio. Retries hard:
# Archive.org's /download/ 302-redirects to per-file storage nodes whose DNS
# can flap, so we retry the whole fetch several times (curl --retry-all-errors
# also re-attempts "could not resolve host"). Returns non-zero only if every
# attempt failed.
download_one() {
    local key="$1" url="$2" attempt
    for attempt in 1 2 3 4 5; do
        if curl -fsSL --ipv4 --retry 4 --retry-delay 2 --retry-all-errors \
                --connect-timeout 30 --max-time 600 -o "$TMPDIR/$key" "$url"; then
            if [ -s "$TMPDIR/$key" ] && \
               { ! command -v file &> /dev/null || file "$TMPDIR/$key" | grep -qiE "audio|mpeg|mp3"; }; then
                return 0
            fi
        fi
        echo "     retry $attempt failed for $key" >&2
        sleep $((attempt * 2))
    done
    return 1
}

echo "Downloading ${#TRACKS[@]} canonical tracks..."
# A single flaky node must not abort the whole seed: collect failures, keep the
# rest, and report. Re-running the script picks up any stragglers (idempotent).
set +e
ok_tracks=()
failed_keys=()
for entry in "${TRACKS[@]}"; do
    key="${entry%%|*}"
    url="${entry#*|}"
    echo "  -> $key"
    if download_one "$key" "$url"; then
        ok_tracks+=("$entry")
    else
        echo "  !! FAILED to download $key — skipping" >&2
        failed_keys+=("$key")
    fi
done
set -e

if [ ${#ok_tracks[@]} -eq 0 ]; then
    echo "ERROR: no tracks downloaded — aborting" >&2
    exit 1
fi
if [ ${#failed_keys[@]} -gt 0 ]; then
    echo "WARNING: ${#failed_keys[@]} track(s) failed and were skipped: ${failed_keys[*]}" >&2
    echo "         Re-run ./scripts/seed-s3.sh to retry them." >&2
fi
# Upload only what we actually downloaded.
TRACKS=("${ok_tracks[@]}")

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
