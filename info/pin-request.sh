curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: alexstorm" \
     -d '{
           "cids": [
             {
               "cid": "QmZEwft3uYvaTp154gdjLyJVSDFjCatNuGTmANGKpvMPZt",
               "labels": {"test1": "value1","test2": "value2"}
             }
           ],
           "rep_min": 2,
           "rep_max": 3
         }' \
     http://localhost:8008/pin