The Organization of the Server Code

The basic packages in EuphoriaDB are structured
hierarchically, in the following order:

* file (Manages OS files as a virtual disk.)
* log (Manages the log.)
* buffer (Manages a buffer pool of pages in memory that acts as a
            cache of disk blocks.)
* tx (Implements transactions at the page level.  Does locking
        and logging.)
* record (Implements fixed-length records inside of pages.)
* metadata (Maintains metadata in the system catalog.)
* query (Implements relational algebra operations. Each operation 
            has a scan class, which can be composed to create a query tree.)
* parse (Implements the parser.)
* plan (Implements a naive planner for SQL statements.)
* jdbc (Implements embedded and network interfaces for JDBC.)
* server (The place where the startup and initialization code live. 
            The class Startup contains the main method.)

The basic server is exceptionally inefficient. The following packages
enable more efficient query processing:

* index (Implements static hash and btree indexes, as well as 
            extensions to the parser and planner to take advantage
            of them.)
* materialize (Implements implementations of the relational 
                operators materialize, sort, groupby, and mergejoin.)
* multibuffer (Implements modifications to the sort and product 
                operators, in order to make optimum use of available
                buffers.)
* opt (Implements a heuristic query optimizer)