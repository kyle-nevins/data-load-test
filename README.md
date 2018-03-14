# data-load-test

Using the PGLoader utility

1. Upload the utility and bind the service instance to the app.
1. In your loader configuration file;
  1. The FROM line should be the name of your data file or the URI to
  the data. If it is a data file, do not include any folder prefix, just the
  name of the file only.
    - FROM 'example_pgloader_COPY_load.sql'
  1. The INTO line should must have the key
  postgresql://username:password@host:port/database. Follow that with a ? and
  the name of the table into which the data goes if needed.
    - INTO postgresql://username:password@host:port/database?pgloader_data_load
1. Command passed is `pgloader -f loaderFile`

Using the PSQL utility
1.

Using the PGRestore Utility
1. Upload the utility and bind the service instance to the app.
1. PGRestore uses the following flags:
  1. --clean --jobs=10 --no-owner --no-privileges
1.  
