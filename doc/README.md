# User Management System Enhancement - Documentation Index

## Overview

This directory contains comprehensive documentation for the User Management System Enhancement project, which migrates from a WeChat-dependent authentication system to a game-account-based user management system.

---

## Project Goals

1. **Remove WeChat Dependencies**: Eliminate reliance on WeChat for user authentication and profile management
2. **Game Account Integration**: Require all users to bind game accounts during registration
3. **Multi-Account Support**: Enable super administrators to manage multiple game accounts with store bindings
4. **Automatic Synchronization**: Auto-activate store sessions and sync game records when accounts are bound
5. **Access Control**: Enforce strict constraints on store administrator permissions

---

## Documentation Structure

### 1. [Implementation Plan](./implementation_plan.md)
**Purpose:** High-level roadmap for the entire project

**Contents:**
- Executive summary
- Phase-by-phase breakdown
- Database schema updates
- Backend implementation details
- Frontend implementation details
- Testing and validation strategies
- Success criteria

**Audience:** Project managers, technical leads, developers

**When to Use:** 
- Planning project timeline
- Understanding overall architecture
- Coordinating team efforts

---

### 2. [Technical Specifications](./technical_specifications.md)
**Purpose:** Detailed technical requirements and code examples

**Contents:**
- Database layer specifications with SQL schemas
- Backend API specifications with Go code examples
- Frontend component specifications with TypeScript/React examples
- Integration specifications
- Security specifications
- Transaction handling examples
- Concurrent processing patterns

**Audience:** Backend developers, frontend developers, database administrators

**When to Use:**
- Implementing specific features
- Understanding data models
- Writing code
- Reviewing technical design

---

### 3. [Migration Guide](./migration_guide.md)
**Purpose:** Step-by-step instructions for executing the migration

**Contents:**
- Pre-migration checklist
- Database migration steps with SQL scripts
- Backend deployment procedures
- Frontend deployment procedures
- Data migration for existing users
- Post-migration cleanup
- Rollback procedures
- Monitoring and alerts
- Troubleshooting guide

**Audience:** DevOps engineers, database administrators, system administrators

**When to Use:**
- Executing the migration
- Planning deployment
- Handling rollbacks
- Troubleshooting issues

---

### 4. [API Documentation](./api_documentation.md)
**Purpose:** Complete API reference for all endpoints

**Contents:**
- Authentication APIs
- User management APIs
- Game account management APIs
- Super admin APIs
- Session management APIs
- Game record APIs
- Error codes reference
- Rate limiting information
- SDK examples (JavaScript, Go)

**Audience:** Frontend developers, API consumers, integration partners

**When to Use:**
- Integrating with the API
- Understanding request/response formats
- Implementing client applications
- Debugging API calls

---

### 5. [Database DDL Scripts](./user_management.ddl)
**Purpose:** SQL scripts for database schema changes

**Contents:**
- Foreign key constraints
- Unique constraints for one-to-one relationships
- Partial unique indexes for soft-deleted records
- Player-dimension record table creation
- Index creation for query optimization

**Audience:** Database administrators, backend developers

**When to Use:**
- Applying database changes
- Understanding schema structure
- Creating test databases

---

### 6. [Full Schema Export Script](./generate_full_schema.sh)
**Purpose:** Bash script to export complete database schema

**Contents:**
- pg_dump command with proper flags
- Configuration variables
- Usage instructions

**Audience:** Database administrators, DevOps engineers

**When to Use:**
- Documenting current schema
- Creating backups
- Setting up development environments

---

## Quick Start Guide

### For Developers

1. **Read First:**
   - [Implementation Plan](./implementation_plan.md) - Understand the big picture
   - [Technical Specifications](./technical_specifications.md) - Get implementation details

2. **During Development:**
   - Reference [API Documentation](./api_documentation.md) for endpoint specs
   - Check [Technical Specifications](./technical_specifications.md) for code examples
   - Review [user_management.ddl](./user_management.ddl) for database schema

3. **Testing:**
   - Follow testing guidelines in [Implementation Plan](./implementation_plan.md) Phase 4
   - Use [API Documentation](./api_documentation.md) for API testing

### For DevOps/DBAs

1. **Read First:**
   - [Migration Guide](./migration_guide.md) - Complete deployment procedure
   - [user_management.ddl](./user_management.ddl) - Database changes

2. **Before Migration:**
   - Complete pre-migration checklist in [Migration Guide](./migration_guide.md)
   - Run [generate_full_schema.sh](./generate_full_schema.sh) to backup current schema
   - Review rollback procedures

3. **During Migration:**
   - Follow [Migration Guide](./migration_guide.md) step-by-step
   - Monitor metrics defined in monitoring section
   - Keep rollback plan ready

4. **After Migration:**
   - Verify data integrity using queries in [Migration Guide](./migration_guide.md)
   - Monitor performance metrics
   - Complete post-migration cleanup

### For Project Managers

1. **Planning:**
   - Review [Implementation Plan](./implementation_plan.md) for timeline and phases
   - Check success criteria and deliverables
   - Understand resource requirements

2. **Tracking:**
   - Use phase breakdown for milestone tracking
   - Monitor success criteria completion
   - Review testing and validation progress

---

## Key Concepts

### Game Account Binding

**What:** Linking a platform user account to a game account

**Why:** 
- Enables automatic profile synchronization
- Provides single source of truth for user identity
- Facilitates game record tracking

**How:**
- During registration: Required field, validated against game API
- After registration: Optional for regular users, multiple for super admins

### Store Session

**What:** Active connection between a game account and a store

**Why:**
- Enables real-time game data monitoring
- Triggers automatic record synchronization
- Manages access control

**How:**
- Auto-activated when super admin binds game account to store
- Can be manually started/stopped
- Monitored by background service

### Player-Dimension Records

**What:** Game records organized by individual players rather than by games

**Why:**
- Efficient player statistics queries
- Better data normalization
- Supports player-centric analytics

**How:**
- One record per player per game
- Links to both battle record and game account
- Indexed for fast queries

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (battle-reusables)              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Registration │  │ Super Admin  │  │ Store Admin  │      │
│  │    Flow      │  │  Dashboard   │  │  Dashboard   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTPS/REST API
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Backend (battle-tiles)                     │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                   Service Layer                       │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐     │   │
│  │  │   Auth     │  │   Game     │  │   Store    │     │   │
│  │  │  Service   │  │  Account   │  │  Service   │     │   │
│  │  └────────────┘  └────────────┘  └────────────┘     │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                  Business Logic Layer                 │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐     │   │
│  │  │   User     │  │   Game     │  │   Sync     │     │   │
│  │  │  UseCase   │  │  Account   │  │  UseCase   │     │   │
│  │  └────────────┘  └────────────┘  └────────────┘     │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                   Data Access Layer                   │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐     │   │
│  │  │   User     │  │   Game     │  │   Store    │     │   │
│  │  │   Repo     │  │  Account   │  │   Repo     │     │   │
│  │  └────────────┘  └────────────┘  └────────────┘     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ SQL
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Database (PostgreSQL)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  basic_user  │  │game_account  │  │game_player   │      │
│  │              │  │              │  │   _record    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │game_account  │  │game_shop     │  │game_session  │      │
│  │   _house     │  │   _admin     │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Background Jobs
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Task Queue (Asynq)                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │         TaskSyncGameRecords                           │   │
│  │  - Fetch battle records from game API                │   │
│  │  - Transform to player-dimension records             │   │
│  │  - Save to database                                  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Database Schema Diagram

```
┌─────────────────┐
│  basic_user     │
│─────────────────│
│ id (PK)         │◄─────────┐
│ username        │          │
│ password        │          │
│ nick_name       │          │
│ avatar          │          │
│ wechat_id       │          │ (deprecated)
└─────────────────┘          │
                             │
                             │ user_id (FK)
                             │
┌─────────────────┐          │
│ game_account    │          │
│─────────────────│          │
│ id (PK)         │──────────┘
│ user_id (FK)    │
│ account         │
│ nickname        │
│ ctrl_account_id │──────┐
└─────────────────┘      │
        │                │
        │ game_account_id│ ctrl_account_id (FK)
        │                │
        ▼                ▼
┌─────────────────┐  ┌─────────────────┐
│game_account     │  │game_ctrl        │
│   _house        │  │   _account      │
│─────────────────│  │─────────────────│
│ id (PK)         │  │ id (PK)         │
│ game_account_id │  │ identifier      │
│ house_gid       │  │ game_user_id    │
│ [UK: account_id]│  │ status          │
└─────────────────┘  └─────────────────┘
        │                    │
        │ house_gid          │ ctrl_account_id
        │                    │
        ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│ game_shop       │  │ game_session    │
│   _admin        │  │─────────────────│
│─────────────────│  │ id (PK)         │
│ id (PK)         │  │ ctrl_account_id │
│ user_id (FK)    │  │ house_gid       │
│ house_gid       │  │ state           │
│ role            │  │ device_ip       │
│ [UK: user_id    │  └─────────────────┘
│  WHERE deleted  │
│  _at IS NULL]   │
└─────────────────┘

┌─────────────────┐
│ game_battle     │
│   _record       │
│─────────────────│
│ id (PK)         │◄─────────┐
│ house_gid       │          │
│ group_id        │          │
│ room_uid        │          │
│ players_json    │          │ battle_record_id (FK)
│ battle_at       │          │
└─────────────────┘          │
                             │
                             │
┌─────────────────┐          │
│ game_player     │          │
│   _record       │          │
│─────────────────│          │
│ id (PK)         │──────────┘
│ battle_record_id│
│ player_gid      │
│ game_account_id │
│ score_delta     │
│ is_winner       │
│ battle_at       │
│ meta_json       │
└─────────────────┘
```

---

## Constraints Summary

| Table | Constraint | Type | Purpose |
|-------|------------|------|---------|
| game_account | fk_game_account_user | Foreign Key | Link to basic_user |
| game_account_house | uk_game_account_house_account | Unique | One store per game account |
| game_shop_admin | uk_game_shop_admin_user_active | Partial Unique Index | One active store per user |
| game_player_record | fk_gpr_battle | Foreign Key | Link to battle record |
| game_player_record | fk_gpr_game_account | Foreign Key | Link to game account |

---

## Glossary

- **Basic User**: Platform user account with username/password authentication
- **Game Account**: User's game credentials and profile
- **Ctrl Account**: System-level game account for automation
- **House/Store**: Game venue or store location
- **Session**: Active connection between game account and store
- **Battle Record**: Single game/match record
- **Player Record**: Individual player's participation in a game
- **Super Admin**: User with full system access and multi-account management
- **Store Admin**: User managing a single store
- **Game User ID**: Unique identifier from game system
- **Player GID**: Game player global identifier

---

## Support and Contact

- **Technical Issues**: Create issue in project repository
- **Documentation Updates**: Submit pull request
- **Questions**: Contact development team

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-11-11 | Initial documentation release |

---

## License

Internal documentation - Proprietary and confidential

