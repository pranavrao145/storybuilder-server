IP_ADDR=$(shell hostname -i | tr -s ' ' | cut -d ' ' -f 1)

default: reset build run

run:
	@echo "Starting StoryBuilder server..."
	docker-compose up -d
	@echo "StoryBuilder server started successfully. When using StoryBuilder CLI, use http://$(IP_ADDR):8080 as your SERVER."

build:
	@echo "Updating local repository..."
	git pull origin master
	@echo "Building StoryBuilder server..."
	docker-compose build
	@echo "StoryBuilder server built successfully."


clean stop reset:
	@echo "Stopping StoryBuilder server."
	docker-compose down --rmi all || true
	@echo "StoryBuilder server stopped and reset successfully."
