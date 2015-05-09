# glitchname
Generates names and checks their availability on Twitter

## generating names
Uses provided name as set of characters and generates all available subsets
of that set.

## checking availability
Checks availability using HTTP status of URL `http://twitter.com/<username>`.
If returns 404, the name is considered available.

Automatically rejects names shorter than 4.

## usage
```
  -name="": base name used for generation
  -sleep=500: sleep in ms between requests
  -verbose=false: output all names, even taken ones
  -workers=4: number of workers querying Twitter
```
