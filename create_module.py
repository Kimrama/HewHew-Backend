import sys
import os

try:
    module_name = sys.argv[1]
except IndexError:
    print("No module name provided!")
    exit(0)

pkg_path = os.path.join(os.getcwd(), "pkg")

def pascal_to_camel(s):
    return s[0].lower() + s[1:]

# create module directory
module_directory = os.path.join(pkg_path, f"{pascal_to_camel(module_name)}")
os.mkdir(module_directory)

# create subdirectories
subdirs = ["controller", "service", "repository", "model"]
for subdir in subdirs:
    os.mkdir(os.path.join(module_directory, subdir))

# file schema =================================================
repository_schema = \
f"""
package repository

type {module_name}Repository interface {{}}
"""
repository_implement_schema = \
f"""
package repository
import ("hewhew-backend/database")
type {module_name}RepositoryImpl struct {{
    db database.Database
}}
func New{module_name}RepositoryImpl(db database.Database) {module_name}Repository {{
return &{module_name}RepositoryImpl{{
    db: db,
}}
}}
"""
service_schema = \
f"""
package service

type {module_name}Service interface {{}}
"""
service_implement_schema = \
f"""
package service

import (
    "hewhew-backend/pkg/{pascal_to_camel(module_name)}/repository"
)

type {module_name}ServiceImpl struct {{
    {module_name}Repository repository.{module_name}Repository
}}

func New{module_name}ServiceImpl({module_name}Repository repository.{module_name}Repository) {module_name}Service {{
    return &{module_name}ServiceImpl{{
        {module_name}Repository: {module_name}Repository,
    }}
}}
"""
controller_schema = \
f"""
package controller

type {module_name}Controller interface {{}}
"""
controller_implement_schema = \
    f"""
package controller

import (
    "hewhew-backend/pkg/{pascal_to_camel(module_name)}/service"
)

type {module_name}ControllerImpl struct {{
    {module_name}Service service.{module_name}Service
}}

func New{module_name}ControllerImpl({module_name}Service service.{module_name}Service) {module_name}Controller {{
    return &{module_name}ControllerImpl{{
        {module_name}Service: {module_name}Service,
    }}
}}
"""
# create repository =======================================================
repository_file_path = os.path.join(module_directory, "repository", f"{pascal_to_camel(module_name)}Repository.go")

with open(repository_file_path, "w") as f:
    f.write(repository_schema)

with open(os.path.join(module_directory, "repository", f"{pascal_to_camel(module_name)}RepositoryImpl.go"), "w") as f:
    f.write(repository_implement_schema)

# create service =======================================================
service_file_path = os.path.join(module_directory, "service", f"{pascal_to_camel(module_name)}Service.go")

with open(service_file_path, "w") as f:
    f.write(service_schema)

with open(os.path.join(module_directory, "service", f"{pascal_to_camel(module_name)}ServiceImpl.go"), "w") as f:
    f.write(service_implement_schema)

# create controller =======================================================
controller_file_path = os.path.join(module_directory, "controller", f"{pascal_to_camel(module_name)}Controller.go")

with open(controller_file_path, "w") as f:
    f.write(controller_schema)

with open(os.path.join(module_directory, "controller", f"{pascal_to_camel(module_name)}ControllerImpl.go"), "w") as f:
    f.write(controller_implement_schema)

# create init router
init_router_schema = \
f"""
package server

import (
	_{pascal_to_camel(module_name)}Controller "hewhew-backend/pkg/{pascal_to_camel(module_name)}/controller"
	_{pascal_to_camel(module_name)}Repository "hewhew-backend/pkg/{pascal_to_camel(module_name)}/repository"
	_{pascal_to_camel(module_name)}Service "hewhew-backend/pkg/{pascal_to_camel(module_name)}/service"
)

func (s *fiberServer) init{module_name}Router() {{
	{pascal_to_camel(module_name)}Repository := _{pascal_to_camel(module_name)}Repository.New{module_name}RepositoryImpl(s.db)
	{pascal_to_camel(module_name)}Service := _{pascal_to_camel(module_name)}Service.New{module_name}ServiceImpl({pascal_to_camel(module_name)}Repository)
	{pascal_to_camel(module_name)}Controller := _{pascal_to_camel(module_name)}Controller.New{module_name}ControllerImpl({pascal_to_camel(module_name)}Service)

	{pascal_to_camel(module_name)}Group := s.app.Group("/v1/{pascal_to_camel(module_name)}")
}}
"""
with open(os.path.join(os.getcwd(), "server", f"init{module_name}Router.go"), "w") as f:
    f.write(init_router_schema)

# server call init router
with open(os.path.join(os.getcwd(), "server", "server.go"), "r") as f:
    server_go_content = f.readlines()
    modified_server_go_content = []
    for line in server_go_content:
        modified_server_go_content.append(line)
        if line.strip() == "// Initialize routes":
            modified_server_go_content.append("	s.init" + module_name + "Router()\n")
            
with open(os.path.join(os.getcwd(), "server", "server.go"), "w") as f:
    f.writelines(modified_server_go_content)
print("create module complete")