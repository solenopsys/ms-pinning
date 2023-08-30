curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: alexstorm" \
     -d '{
           "pins": [
             {
               "cid": "QmQb25BukhubpfTrpSjtunk2pWEgURqbrMD6NU1vFnanPq",
               "labels": {"site": "robotization.vc"}
             }
           ],
           "rep_min": 2,
           "rep_max": 3
         }' \
     http://pinning.solenopsys.org/pin
