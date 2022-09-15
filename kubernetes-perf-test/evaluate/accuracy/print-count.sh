fileName="100Users.jtl"
sampleCount=3

for i in `seq 1 $sampleCount`
do
    echo "filtered_${i}_${fileName}"
    printf "200,OK: %s\n" $(grep "200,OK" -o "filtered_${i}_${fileName}" | wc -l)
    printf "500,Internal Server Error: %s\n" $(grep "500,Internal Server Error" -o "filtered_${i}_${fileName}" | wc -l)
    printf "429,Too Many Requests: %s\n" $(grep "429,Too Many Requests" -o "filtered_${i}_${fileName}" | wc -l)
    printf "Non 200,OK: %s\n" $(grep -v "200,OK" -o "filtered_${i}_${fileName}" | wc -l)
    echo ""
done

