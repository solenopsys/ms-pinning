curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: alexstorm" \
     -d '{
           "cids": [
             "QmX374s4EKeF57rMAkyJmgxqmPHUcBNET4eHMqpi7G5dGj"
           ],
           "rep_min": 2,
           "rep_max": 3
         }' \
     http://pinning.solenopsys.org/pin