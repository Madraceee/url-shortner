# URL Shortner 
Built using GO

---
The URL gets converted to base64, then the first 7 digits are used as hash. 

---
| API | Description |
| ----------- | ----------- |
| /health | Returns Server Status |
| /v1/shortURL | Pass a JSON object with "longURL" as key and link as parameter. The shortURL is sent as response |
| /v1/{shortURL} | Pass the shortURL and page gets redirected to the initial Site |

---
### Future Updates
    1. Check the urls using regex instead of matching
    2. Using redis to improve speed and reduce load on server
