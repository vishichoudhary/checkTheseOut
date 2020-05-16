### Directory Structure
1. `common` contains data structures which are universal throughout the service.
2. `config` contains the configuration for the server. 
3. `pdfLogs` contains the generated pdfs, the format of pdf is `By_uuid of user_at_timestamp` which represent pdf belong to which user and when it is created.


### How it works
UserID is attached in body of the request as adding authentication and identification is beyond the scope of this assignment. If a valid request is made to the service, service will start a stopwatch of 5 minutes. If during this duration, another request of the same user is encountered, the service will update the JSON data and restart the stopwatch with 5 minutes. After the completion of 5 minutes a pdf with name `By_uuid of user_at_timestamp` is created under `pdfLogs/` folder. 

### Request Example :
URL : `<hostname>/v1/generatePdf`

Http Method : `POST`

Body : 
```
{
    "UserID": "4507e41e-9773-11ea-bb37-0242ac130002",
    "Questions": [
        "When is the last time you experienced nostalgia?",
        "What's the scariest dream you've ever had?",
        "What's the weirdest thought you've ever had?",
        "What's the first thing that comes to mind for fidget?"
    ]
}
```

### Failure Response

Http Response Code : != 200

Body:
```{}```