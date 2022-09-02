import csv
csvFile = csv.reader(open(r"sample.jtl"),delimiter=',')
filtered = filter(lambda row: row[0] == 'timeStamp' or 1661511913000 <= int(row[0]) < 1661511914000, csvFile)
csv.writer(open(r"filtered.jtl",'w'),delimiter=',').writerows(filtered)
