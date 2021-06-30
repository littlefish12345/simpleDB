# simpleDB
A file based database library

# Notice
This database library is only for my personal usage. It is possible that I will not update it anymore or add support for other types of data. It can only be use as a map in file and only support int64 for key and int64 for value.

# API
simpleDB.CreateDatabase(path string) (error) Create a database.  
simpleDB.OpenDatabase(path string) (error) Open a database.  
simpleDB.CloseDatabase() Close a database.  
simpleDB.WriteDatabase(key int64, value int64) (error) Write or change a value to an opened database.  
simpleDB.ReadDatabase(key int64) (int64, error) Read a value from an opened database.  

# Errors
simpleDB.DBDamaged This means that the database is damaged or it is not a simpleDB database.  
simpleDB.DBNotOpened This means that you haven't open a database.  
simpleDB.DBKeyNotFound This means that the key you have requested not found.