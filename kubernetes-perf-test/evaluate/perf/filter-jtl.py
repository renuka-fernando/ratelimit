import csv

# for perf test do not need NodeJS File
firstTime = 1662516059733
fileName = "1000users_100ms_NewRoutes.jtl"

startTime = firstTime + (5 * 60 * 1000) # add 5 min
endTime = firstTime + (15 * 60 * 1000) # add 15 min

csvFile = csv.reader(open(fileName),delimiter=',')
filtered = filter(lambda row: row[0] == 'timeStamp' or startTime <= int(row[0]) < endTime, csvFile)
csv.writer(open("filtered_1_" + fileName,'w'),delimiter=',').writerows(filtered)
