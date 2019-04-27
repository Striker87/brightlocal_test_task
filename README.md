# brightlocal_test_task

Task Create an application providing HTTP based API with JSON message format that provides access to key-value storage. Storage should be in memory, up to 1024 keys must be supported. 
 
Methods: 
● SET: sets a value to a key 
● GET: returns value by key 
● DELETE: removes value 
● EXISTS: tells if key exists 
 
Key is a unique string, cannot be empty, up to 16 bytes long. Value is an arbitrary string, up to 512 bytes long. 
 
Do not rely on specific HTTP properties (headers, URIs, response codes, etc.). The protocol should be portable enough to easily switch to TCP or anything else. 
 
The code should be covered with tests where appropriate. Code should be uploaded to GitHub with multiple staged commits. 
