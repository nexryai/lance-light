## Bug code
A bug code is a uuid that is displayed when an unexpected fatal error occurs within the program that is different from a normal error and makes it impossible to continue processing.

#### `ea14b6c8-aba2-4e2e-8fda-faa448e5771d`
Processing cannot continue because the timestamp gotten from the database was invalid.

- The database may be corrupted. Resetting the database will likely resolve the problem, but this will result in the loss of all previous logs.
- If you encounter this error, please open an issue and contact the developer first. The probability of this error being caused by your error is very low.


#### `8a04693b-9a36-422b-81b6-2270ad8e357b`
This error occurs when fetching a list of IP addresses of Cloudflare servers.  
The request was successful, but an invalid IP address was detected.

 - It is possible that there is a problem with Cloudflare's system, but since the request itself was successful, this is highly unlikely.
 - This is probably caused by a bug in LanceLight; please open an issue and contact the developer.

#### `fea1507a-6eb7-40d4-a499-1f70ac6fd580`
MkAllowPort()関数内にバグがあります。このバグコードはユーザー側のミスでは発生しないはずです。開発者に連絡してください。