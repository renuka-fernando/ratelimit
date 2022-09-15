// in bash find first timestamp using `head <file_name>`

const firstTime = 1663051704591;
const roundToMs = 60000; // 1 min
const addMs = 120000 - 10000; // 2 min - 10 seconds (i.e. 1 min and 50 seconds)

console.log("Start Time:        " + new Date(firstTime))

let startFilterTime = new Date((Math.round(firstTime / roundToMs) * roundToMs) + addMs)
console.log("Start Filter Time: " + startFilterTime)
console.log("Start Filter Time: " + startFilterTime.getTime())
