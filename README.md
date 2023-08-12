# file-clickhouse-exporter
file-clickhouse-exporter is a text file parser and imports the text content into clickhouse.file-clickhouse-exporter is inspired by `awk`.   
```
hello  world  my name  is   tyltr
$1     $2     $3       $4    $4
```

## How to build from sources

### 1. Clone source code
```
git clone git@github.com:tylitianrui/file-clickhouse-exporter.git
```
### 2. Install dependencies
```
cd file-clickhouse-exporter
make deps
```
### 3. Build binary file
```
make build
```

## How it works

### 1. A Simple Example

#### ClickHouse table schema example
This is table schema for our example.
```
CREATE TABLE
    game_score(
        id UUID NOT NULL COMMENT 'Primary KEY',
        user String COMMENT 'user',
        age Int32 COMMENT 'age',
        gender String COMMENT 'user',
        event  String COMMENT 'event',
        score  Float   COMMENT 'score',
        game_start DateTime COMMENT 'game_start',
        game_end DateTime COMMENT 'game_end',
        game String COMMENT 'game',
    ) ENGINE = MergeTree()
ORDER BY (id) PRIMARY KEY(id) COMMENT '';
```

#### text 
**text format**
```
firstName:xxxx   familyName:xxx  <age>  <gender>  <sport_event>  <score> [game_start:xxxx]  [game_end:xxxx] 
```   
 
**text example**    
[text file](./example/usage/data)  
```
firstName:tyltr   familyName:li  28  male  Marathon  210.21 [game_start:2023-07-17T00:02:00+00:00]  [game_end:2023-07-17T00:07:00+00:00]      
firstName:limimg   familyName:jing  34  male  marathon  188.2 [game_start:2023-07-17T00:02:00+00:00]  [game_end:2023-07-17T00:07:00+00:00]         
```
**paser config**  
[config](./example/usage/config.yaml)
```
  columns:
    id: dynamic.id
    user: aggregation.user
    age: $3(int64)
    gender: $4
    event: $5
    score: $6(float64)
    game_start: aggregation.gameStart(time)
    game_end: aggregation.gameEnd(time)
    game: static.game
  Preprocessing:
    aggregation:
      user: $1[10:]+" "+$2[11:]    
      gameStart: $7[12:-1]
      gameEnd: $8[10:-1]
    static:
      game: olympic
    dynamic:
      id: gen_uuid()
```
**run**
```
file-clickhouse-exporter  run
```

## Preprocessing
there are three text preprocessor: `aggregation`  `static` and `dynamic` .
### preprocessor aggregation
aggregation is designed for text join or text split. 

#### **text split**   
usage 
```
$<indexOfText>[from_idx:to_idx]
```
- from_idx   int  
- to_idx     int


```
firstName:tyltr   familyName:li  28  male  Marathon  210.21 [game_start:2023-07-17T00:02:00+00:00]  [game_end:2023-07-17T00:07:00+00:00] 
```
- extract firstName in text :  `$1[:]`    ->  `tyltr`
- extract familyName in text :  `$1[11:]`    ->  `li`
- extract game_start in text :  `$7[12:-1]`    ->  `2023-07-17T00:02:00+00:00`


#### **text join**   

usage 
```
$<indexOfText>[from_idx:to_idx] + "string you want" + $<indexOfText>[from_idx:to_idx] 
```
