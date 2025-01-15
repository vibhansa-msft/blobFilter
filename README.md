# BlobFilter 
Filter module to filter your Azure Storage Blobs based on various fields and patterns

## About
BlobFilter is an open source project developed to filter your blobs based on various conditions. This library can be integrated in any module or application that lists the Azure Storage Blobs using ListBlob API. Once the list of blobs is retrieved you can supply the details of blob to this filter library and check blob properties matches your set pattern or not.

## Filter Types
BlobFilter allows user to filter blobs based on
- Access Tier
- File Extension
- Last modified date
- Blob Name
- Size
- Blob Tag


## Filter Details
* Access Tier
   - Supported operations are "=" and "!="
   - Tier value can be supplied as string
   - e.g. ```tier=hot```

* File Extension
   - Supported operations are "=" and "!="
   - Extension can be supplied as string. Do not include "." in the filter
   - e.g. ```format=pdf```

* Last Modified Date
   - Supported operations are "<=", ">=", "<", ">" and "="
   - Date shall be provided in RFC1123 Format e.g. "Mon, 24 Jan 1982 13:00:00 IST"
   - e.g. ```modtime>Mon, 24 Jan 1982 13:00:00 IST```

* Blob Name
   - Supported operations are "=" and "!="
   - Name shall be a valid regex expression
   - e.g. ```name=^mine[0-1]\\d{3}.*```

* Size
   - Supported operations are "<=", ">=", "!=", "<", ">" and "="
   - Size shall be provided in bytes
   - e.g. ```size > 1000```

* Blob Tag
   - Supported operations are "=" and "!="
   - Tag shall be provided in key:value format e.g. "key1:val2"
   - e.g. ```tag=key1:val1```


## Complex filter expression
- Complex filters can be provided using logical AND and OR operations
- For AND use "&&" and for OR use "||"
- Brackets are not supported as of now
- Filters are divided based on OR operation first and each part is considered a sub-filter jonied with AND
- One sub filter can have multiple conditions joined using AND
- Example
    - ```size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*```
    - This filter has 4 sub-filters
    - ```size > 1000 && tag=key1:val1```, ``` size > 2000 && tag=key2:val2 ```, ```tier=hot``` and ```name=^mine[0-1]\\d{3}.*```
    - As all subfilters are joined using OR first filter that matches with provided properties will terminate any further filtering and it will be considered a hit
    - Each sub-filter may have multiple conditions joined using AND. First condition that does not match will terminate further filtering of that sub filter and declare the result as miss.


## Sample Code (Sync filter)
```
	bf := BlobFilter{}
	bf.Configure("size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*")

   attr := BlobAttr{Size: 1500, Tags: map[string]string{"key1": "val1"}, Tier: "cold", Name: "nine1982.doc"}
   result := bf.IsAcceptable(&attr)
   if result {
      fmt.Println("Blob Matches")
   } else {
      fmt.Println("Blob Does not match")
   }
	
```


## Async filters
- Application that needs to filter large number of objects can use parallel filters.
- By using ```EnableAsyncFilter``` method application can set how many go routines shall BlobFilter run to filter files in parallel
- Once concurrency is set, application can enqueue items using ```AddItem``` API
- This api accepts a key (unique for each item) and the attributes of the Blob
- Results can be retreived back using ```NextResult``` API, which returns back the key and the result
- Application is free to co-relate the key to blobs in anyway it wants
- Once application has processed all results it can stop concurrent filters using ```TerminateAsyncFilter```
- Once concurrency is terminated, filters will no longer work. If application wishes to resume it has to turn the concurrency back on.
- While concurrent filters are running, application can still use ```IsAcceptable``` API which is synchronous call.

## Sample Code (Async filter)
```
	bf := BlobFilter{}
	bf.Configure("size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*")

   attr := BlobAttr{Size: 1500, Tags: map[string]string{"key1": "val1"}, Tier: "cold", Name: "nine1982.doc"}
   bf.AddItem("key1", &attr)


   key, result := bf.NextResult()	
   if result {
      fmt.Println("Blob %s Matches", key)
   } else {
      fmt.Println("Blob %s Does not match", key)
   }
```












