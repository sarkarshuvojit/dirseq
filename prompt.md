Create a golang cli tool which will

1. Check if config dir present at $HOME/.config/spd
1.1. if Not; Create it
2. Check if mem.db is present in spd dir
2.1 if not create a new sql db called mem.db
3. create a table containing columns 
    abs_path VARCHAR NOT NULL
    last_seq INT
4. Check in that table if current abs path is already present
4.1 if present pick it's last_seq else default to 0
5. Now increment it with one
6. print the value 
7. update db with the new last_seq
