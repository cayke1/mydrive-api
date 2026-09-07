# Database Seed

Quick sample data for development and testing.

## What's included

### Users (3 sample users)
- **john@example.com** - ID: `550e8400-e29b-41d4-a716-446655440001`
- **jane@example.com** - ID: `550e8400-e29b-41d4-a716-446655440002`
- **alice@example.com** - ID: `550e8400-e29b-41d4-a716-446655440003`

### Folders
- Root "My Drive" folder for each user
- Subfolders:
  - John: Documents, Photos, Projects
  - Jane: Work
  - Alice: (only root)

### Files
- Sample documents, images, and spreadsheets
- Distributed across folders with realistic metadata (size, mime type, checksum)

## How to load seed data

### First time setup
```bash
make migrate-up  # Run migrations first
make seed        # Load sample data
```

### Reset and reload
```bash
make seed-reset  # Clears all data and reloads
```

### Manual SQL execution
```bash
psql -U user -d mydrive -f infrastructure/postgres/seed.sql
```

## Testing with the API

After seeding, you can use these UUIDs to test:

```bash
# List files for John (owner_id)
curl http://localhost:8080/files?owner_id=550e8400-e29b-41d4-a716-446655440001

# Get a specific folder
curl http://localhost:8080/folders/650e8400-e29b-41d4-a716-446655440002
```

## Modify seed data

Edit `seed.sql` to add/remove users, folders, or files. Keep UUIDs consistent for foreign key references.

Example UUIDs pattern:
- Users: `550e8400-e29b-41d4-a716-446655440XXX`
- Folders: `650e8400-e29b-41d4-a716-446655440XXX`
- Files: `750e8400-e29b-41d4-a716-446655440XXX`
