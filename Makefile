.PHONY: dev-start dev-stop logs-follow sqlc-generate clean help

# Development commands
dev-start:
	docker compose watch

dev-stop:
	docker compose down

logs-follow:
	docker compose logs --follow acacia frontend

# Database commands
sqlc-generate:
	$(MAKE) -C services/acacia sqlc-generate

# General cleanup
clean:
	docker compose down -v
	docker system prune -f

# Help command
help:
	@echo "Available commands:"
	@echo "  dev-start      - Start development environment with docker compose watch"
	@echo "  dev-stop       - Stop development environment"
	@echo "  logs-follow    - Follow logs for acacia and frontend services"
	@echo "  sqlc-generate  - Generate Go code from SQL queries"
	@echo "  clean          - Clean all artifacts and stop containers"
	@echo "  help           - Show this help message"
