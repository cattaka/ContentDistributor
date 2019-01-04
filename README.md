# ContentDistributor

ContentDistributor distributes contents(such as pdf, epub and zip).
It generates URLs for download them.


## How to deploy
- Prepare [Google Cloud SDK](https://cloud.google.com/sdk/)
- Prepare Firebase Project on [Firebase Cloud](https://console.firebase.google.com/)
- Put 2 files, those are from Firebase Console
  - app/firebaseConfig.json
  - app/serviceAccountKey.json
- Test and check with dev_appserver.py
  - $ cd app
  - $ ./dev_appserver.py app.yaml
- Deploy to Google App Engine with gcloud command


## How to use
- Open admin page: https://<<hostname>>/admin/
- Create Distribution
  - Fill bsic info
  - Upload Cover image
  - Upload Distribution Files
- Generated Codes
- Open Codes screen
  - Download json
- Open util from Codes screen


## License

```
Copyright 2019 Takao Sumitomo

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
