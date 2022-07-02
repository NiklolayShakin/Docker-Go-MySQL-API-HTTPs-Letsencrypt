# Go Docker MySQL API HTTPs
Docker based api to get person info by their phone number.

Scheme: 
[                 Docker                    ]
[ DB-container ]<------->[  API-Container   ]
        ||                      ||
[Volume with DB]    [Volume with app sources]

Rename conf.go.examle --> conf.go and customize the settings if needed.
Deploy the project with "docker-compose up -d" command.
The project imports MySQL-dump (if exist any *.sql files at ./database directory), MySQL stores its data on a volume to keep them safe after container stops/restarts.
By default web-server will listen 8080 HTTP, but you can enable HTTPs-mode in conf.go. It will be asking for a certificate from Let's Encrypt any time you run the server (or making changes in the code due to CompileDemon). Be careful because Let's Encrypt has limits for request from one IP. 
The api image of the project contains CompileDemon, this tool will recompile and restart your Go-app every time you've made any changes in source code.
