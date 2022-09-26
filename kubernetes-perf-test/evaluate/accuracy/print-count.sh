set -e
java -jar jtl-splitter-0.4.6-SNAPSHOT.jar -t 5 -s -p -f "${1}Users-API-level.jtl"

cat "${1}Users-API-level-measurement-summary.json"

python3 filter-jtl.py "${1}Users-API-level-measurement.jtl"
fileName="${1}Users-API-level-measurement.jtl"

sampleCount=3

for i in `seq 1 $sampleCount`
do
    echo "filtered_${i}_${fileName}"
    printf "200,OK: %s\n" $(grep "200,OK" -o "filtered_${i}_${fileName}" | wc -l)
    printf "500,Internal Server Error: %s\n" $(grep "500,Internal Server Error" -o "filtered_${i}_${fileName}" | wc -l)
    printf "429,Too Many Requests: %s\n" $(grep "429,Too Many Requests" -o "filtered_${i}_${fileName}" | wc -l)
    printf "Non 200,OK: %s\n" $(grep -v "200,OK" "filtered_${i}_${fileName}" | wc -l)
    echo ""
done
