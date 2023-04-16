curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <access_token>" \
     -d '{
           "cids": [
             "QmZEwft3uYvaTp154gdjLyJVSDFjCatNuGTmANGKpvMPZt"
           ],
           "rep_min": 2,
           "rep_max": 3
         }' \
     http://localhost:32531/pin