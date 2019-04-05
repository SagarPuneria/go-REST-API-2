# Project details:

A REST API implmentation using golang.

Create typicode database using MySQL.

Using fake REST API ``https://jsonplaceholder.typicode.com/albums`` get array of JSON response as shown below:
```sh
[
  {
    "userId": 1,
    "id": 1,
    "title": "quidem molestiae enim"
  },
  {
    "userId": 1,
    "id": 2,
    "title": "sunt qui excepturi placeat culpa"
  },
  {
    "userId": 1,
    "id": 3,
    "title": "omnis laborum odio"
  }
]
```
Insert above date into album table.

Every album id has set photos. To get set of photos by using above album id frame another URL like ``https://jsonplaceholder.typicode.com/photos?albumId=1`` to get array of JSON response for every album id as shown below:
```sh
[
  {
    "albumId": 1,
    "id": 1,
    "title": "accusamus beatae ad facilis cum similique qui sunt",
    "url": "https://via.placeholder.com/600/92c952",
    "thumbnailUrl": "https://via.placeholder.com/150/92c952"
  },
  {
    "albumId": 1,
    "id": 2,
    "title": "reprehenderit est deserunt velit ipsam",
    "url": "https://via.placeholder.com/600/771796",
    "thumbnailUrl": "https://via.placeholder.com/150/771796"
  },
  {
    "albumId": 1,
    "id": 3,
    "title": "officia porro iure quia iusto qui ipsa ut modi",
    "url": "https://via.placeholder.com/600/24f355",
    "thumbnailUrl": "https://via.placeholder.com/150/24f355"
  }
]
```
Insert above date into photo table.

Create a generic GET API endpoint to get specific the records from respective table. Specific table should be mentioned in Query Params ``key=value`` pair like type=album and id=3 or type=photo and id=1 which will give json response body.
```sh
 GET http://localhost:8080/search?type=album&id=3
 GET http://localhost:8080/search?type=photo&id=1
```
