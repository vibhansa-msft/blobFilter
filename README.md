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

* File Extension
   - Supported operations are "=" and "!="
   - Extension can be supplied as string. Do not include "." in the filter

* Last Modified Date
   - Supported operations are "<=", ">=", "<", ">" and "="
   - Date shall be provided in RFC1123 Format e.g. "Mon, 24 Jan 1982 13:00:00 IST"

* Blob Name
   - Supported operations are "=" and "!="
   - Name shall be a valid regex expression

* Size
   - Supported operations are "<=", ">=", "!=", "<", ">" and "="
   - Size shall be provided in bytes

* Blob Tag
   - Supported operations are "=" and "!="
   - Tag shall be provided in key:value format e.g. "key1:val2"


## Complex filter expression
- Complex filters can be provided using logical AND and OR operations
- For AND use "&&" and for OR use "||"
- Brackets are not supported as of now
- Filters are divided based on OR operation and each part is considered a sub-filter
- One sub filter can have multiple conditions joined using AND
- Example
    - ```size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*```
    - This filter has 4 sub-filters
    - ```size > 1000 && tag=key1:val1``` ``` size > 2000 && tag=key2:val2 ``` ```tier=hot``` and ```name=^mine[0-1]\\d{3}.*```
    - As all subfilters are joined using OR first filter that matches with provided properties will terminate any further filtering and it will be considered a hit
    - Each sub-filter may have multiple conditions joined using AND. First condition that does not match will terminate further filtering of that sub filter and declare the result as miss.














