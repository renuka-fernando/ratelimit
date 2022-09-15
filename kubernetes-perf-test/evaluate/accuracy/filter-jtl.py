import csv

fileName = "100Users.jtl"
startTime = 1663051790000
duration = 60000
sampleCount = 3

startTime2 = startTime + duration
startTime3 = startTime2 + duration
startTime4 = startTime3 + duration

for i in range(1, sampleCount+1):
    endTime = startTime + duration
    csvFile = csv.reader(open(fileName),delimiter=',')
    filtered = filter(lambda row: row[0] == 'timeStamp' or startTime <= int(row[0]) < endTime, csvFile)
    csv.writer(open("filtered_" + str(i) + "_" + fileName,'w'),delimiter=',').writerows(filtered)
    startTime = endTime
