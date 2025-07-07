package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"testing"
)

func TestStripMargin_WithSimpleSelectStatement_ReturnsValidSQL(t *testing.T) {
	// Given
	input := `
	|SELECT id, name, email
	|FROM users
	|WHERE active = true
	|ORDER BY name
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "SELECT id, name, email\nFROM users\nWHERE active = true\nORDER BY name"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithComplexJoinQuery_PreservesQueryStructure(t *testing.T) {
	// Given
	input := `
	|SELECT u.name, p.title, c.name as category
	|FROM users u
	|JOIN posts p ON u.id = p.user_id
	|LEFT JOIN categories c ON p.category_id = c.id
	|WHERE u.active = true
	|  AND p.published_at IS NOT NULL
	|ORDER BY p.created_at DESC
	|LIMIT 10
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "SELECT u.name, p.title, c.name as category\nFROM users u\nJOIN posts p ON u.id = p.user_id\nLEFT JOIN categories c ON p.category_id = c.id\nWHERE u.active = true\n  AND p.published_at IS NOT NULL\nORDER BY p.created_at DESC\nLIMIT 10"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithInsertStatement_ReturnsValidInsert(t *testing.T) {
	// Given
	input := `
	|INSERT INTO users (name, email, created_at)
	|VALUES 
	|  ('John Doe', 'john@example.com', NOW()),
	|  ('Jane Smith', 'jane@example.com', NOW()),
	|  ('Bob Johnson', 'bob@example.com', NOW())
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "INSERT INTO users (name, email, created_at)\nVALUES \n  ('John Doe', 'john@example.com', NOW()),\n  ('Jane Smith', 'jane@example.com', NOW()),\n  ('Bob Johnson', 'bob@example.com', NOW())"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithUpdateStatement_PreservesUpdateLogic(t *testing.T) {
	// Given
	input := `
	|UPDATE users 
	|SET 
	|  name = 'Updated Name',
	|  email = 'updated@example.com',
	|  updated_at = NOW()
	|WHERE id = 123
	|  AND active = true
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "UPDATE users \nSET \n  name = 'Updated Name',\n  email = 'updated@example.com',\n  updated_at = NOW()\nWHERE id = 123\n  AND active = true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithDeleteStatement_ReturnsValidDelete(t *testing.T) {
	// Given
	input := `
	|DELETE FROM posts
	|WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR)
	|  AND status = 'draft'
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "DELETE FROM posts\nWHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR)\n  AND status = 'draft'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithCreateTableStatement_PreservesTableStructure(t *testing.T) {
	// Given
	input := `
	|CREATE TABLE users (
	|  id BIGINT PRIMARY KEY AUTO_INCREMENT,
	|  name VARCHAR(255) NOT NULL,
	|  email VARCHAR(255) UNIQUE NOT NULL,
	|  password_hash VARCHAR(255) NOT NULL,
	|  active BOOLEAN DEFAULT true,
	|  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	|  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	|)
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "CREATE TABLE users (\n  id BIGINT PRIMARY KEY AUTO_INCREMENT,\n  name VARCHAR(255) NOT NULL,\n  email VARCHAR(255) UNIQUE NOT NULL,\n  password_hash VARCHAR(255) NOT NULL,\n  active BOOLEAN DEFAULT true,\n  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n)"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithSubquery_PreservesNestedQuery(t *testing.T) {
	// Given
	input := `
	|SELECT u.name, u.email,
	|  (SELECT COUNT(*) FROM posts p WHERE p.user_id = u.id) as post_count
	|FROM users u
	|WHERE u.id IN (
	|  SELECT DISTINCT user_id 
	|  FROM posts 
	|  WHERE published_at > DATE_SUB(NOW(), INTERVAL 30 DAY)
	|)
	|ORDER BY post_count DESC
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "SELECT u.name, u.email,\n  (SELECT COUNT(*) FROM posts p WHERE p.user_id = u.id) as post_count\nFROM users u\nWHERE u.id IN (\n  SELECT DISTINCT user_id \n  FROM posts \n  WHERE published_at > DATE_SUB(NOW(), INTERVAL 30 DAY)\n)\nORDER BY post_count DESC"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithComments_PreservesComments(t *testing.T) {
	// Given
	input := `
	|-- This query finds active users
	|SELECT id, name, email
	|FROM users
	|WHERE active = true  -- Only active users
	|  AND created_at > '2024-01-01'  /* Users from this year */
	|ORDER BY name
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "-- This query finds active users\nSELECT id, name, email\nFROM users\nWHERE active = true  -- Only active users\n  AND created_at > '2024-01-01'  /* Users from this year */\nORDER BY name"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithStoredProcedure_PreservesStoredProcedure(t *testing.T) {
	// Given
	input := `
	|DELIMITER //
	|CREATE PROCEDURE GetUserStats(IN user_id INT)
	|BEGIN
	|  SELECT 
	|    COUNT(*) as total_posts,
	|    MAX(created_at) as last_post_date
	|  FROM posts 
	|  WHERE posts.user_id = user_id;
	|END //
	|DELIMITER ;
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "DELIMITER //\nCREATE PROCEDURE GetUserStats(IN user_id INT)\nBEGIN\n  SELECT \n    COUNT(*) as total_posts,\n    MAX(created_at) as last_post_date\n  FROM posts \n  WHERE posts.user_id = user_id;\nEND //\nDELIMITER ;"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithCTEQuery_PreservesWithClause(t *testing.T) {
	// Given
	input := `
	|WITH active_users AS (
	|  SELECT id, name, email
	|  FROM users
	|  WHERE active = true
	|),
	|user_posts AS (
	|  SELECT user_id, COUNT(*) as post_count
	|  FROM posts
	|  GROUP BY user_id
	|)
	|SELECT au.name, au.email, COALESCE(up.post_count, 0) as posts
	|FROM active_users au
	|LEFT JOIN user_posts up ON au.id = up.user_id
	|ORDER BY posts DESC
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "WITH active_users AS (\n  SELECT id, name, email\n  FROM users\n  WHERE active = true\n),\nuser_posts AS (\n  SELECT user_id, COUNT(*) as post_count\n  FROM posts\n  GROUP BY user_id\n)\nSELECT au.name, au.email, COALESCE(up.post_count, 0) as posts\nFROM active_users au\nLEFT JOIN user_posts up ON au.id = up.user_id\nORDER BY posts DESC"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithStringLiteralsContainingPipes_PreservesLiterals(t *testing.T) {
	// Given
	input := `
	|SELECT 
	|  id,
	|  CONCAT(first_name, '|', last_name) as full_name,
	|  'Status: Active|Verified' as status_info
	|FROM users
	|WHERE description LIKE '%pipe|symbol%'
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "SELECT \n  id,\n  CONCAT(first_name, '|', last_name) as full_name,\n  'Status: Active|Verified' as status_info\nFROM users\nWHERE description LIKE '%pipe|symbol%'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithWindowFunctions_PreservesWindowSyntax(t *testing.T) {
	// Given
	input := `
	|SELECT 
	|  name,
	|  salary,
	|  ROW_NUMBER() OVER (ORDER BY salary DESC) as rank,
	|  AVG(salary) OVER (PARTITION BY department) as dept_avg
	|FROM employees
	|ORDER BY rank
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "SELECT \n  name,\n  salary,\n  ROW_NUMBER() OVER (ORDER BY salary DESC) as rank,\n  AVG(salary) OVER (PARTITION BY department) as dept_avg\nFROM employees\nORDER BY rank"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

// StripColumn tests for SQL

func TestStripColumn_WithSimpleSelectStatement_ReturnsValidSQL(t *testing.T) {
	// Given
	input := `
	|SELECT id, name, email|
	|FROM users|
	|WHERE active = true|
	|ORDER BY name|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT id, name, email\nFROM users\nWHERE active = true\nORDER BY name"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithComplexJoinQuery_PreservesQueryStructure(t *testing.T) {
	// Given
	input := `
	|SELECT u.name, p.title, c.name as category|
	|FROM users u|
	|JOIN posts p ON u.id = p.user_id|
	|LEFT JOIN categories c ON p.category_id = c.id|
	|WHERE u.active = true|
	|  AND p.published_at IS NOT NULL|
	|ORDER BY p.created_at DESC|
	|LIMIT 10|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT u.name, p.title, c.name as category\nFROM users u\nJOIN posts p ON u.id = p.user_id\nLEFT JOIN categories c ON p.category_id = c.id\nWHERE u.active = true\n  AND p.published_at IS NOT NULL\nORDER BY p.created_at DESC\nLIMIT 10"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithInsertStatement_ReturnsValidInsert(t *testing.T) {
	// Given
	input := `
	|INSERT INTO users (name, email, created_at)|
	|VALUES |
	|  ('John Doe', 'john@example.com', NOW()),|
	|  ('Jane Smith', 'jane@example.com', NOW()),|
	|  ('Bob Johnson', 'bob@example.com', NOW())|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "INSERT INTO users (name, email, created_at)\nVALUES \n  ('John Doe', 'john@example.com', NOW()),\n  ('Jane Smith', 'jane@example.com', NOW()),\n  ('Bob Johnson', 'bob@example.com', NOW())"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithUpdateStatement_PreservesUpdateLogic(t *testing.T) {
	// Given
	input := `
	|UPDATE users |
	|SET |
	|  name = 'Updated Name',|
	|  email = 'updated@example.com',|
	|  updated_at = NOW()|
	|WHERE id = 123|
	|  AND active = true|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "UPDATE users \nSET \n  name = 'Updated Name',\n  email = 'updated@example.com',\n  updated_at = NOW()\nWHERE id = 123\n  AND active = true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithDeleteStatement_ReturnsValidDelete(t *testing.T) {
	// Given
	input := `
	|DELETE FROM posts|
	|WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR)|
	|  AND status = 'draft'|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "DELETE FROM posts\nWHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR)\n  AND status = 'draft'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithCreateTableStatement_PreservesTableStructure(t *testing.T) {
	// Given
	input := `
	|CREATE TABLE users (|
	|  id BIGINT PRIMARY KEY AUTO_INCREMENT,|
	|  name VARCHAR(255) NOT NULL,|
	|  email VARCHAR(255) UNIQUE NOT NULL,|
	|  password_hash VARCHAR(255) NOT NULL,|
	|  active BOOLEAN DEFAULT true,|
	|  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,|
	|  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP|
	|)|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "CREATE TABLE users (\n  id BIGINT PRIMARY KEY AUTO_INCREMENT,\n  name VARCHAR(255) NOT NULL,\n  email VARCHAR(255) UNIQUE NOT NULL,\n  password_hash VARCHAR(255) NOT NULL,\n  active BOOLEAN DEFAULT true,\n  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,\n  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n)"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithSubquery_PreservesNestedQuery(t *testing.T) {
	// Given
	input := `
	|SELECT u.name, u.email,|
	|  (SELECT COUNT(*) FROM posts p WHERE p.user_id = u.id) as post_count|
	|FROM users u|
	|WHERE u.id IN (|
	|  SELECT DISTINCT user_id |
	|  FROM posts |
	|  WHERE published_at > DATE_SUB(NOW(), INTERVAL 30 DAY)|
	|)|
	|ORDER BY post_count DESC|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT u.name, u.email,\n  (SELECT COUNT(*) FROM posts p WHERE p.user_id = u.id) as post_count\nFROM users u\nWHERE u.id IN (\n  SELECT DISTINCT user_id \n  FROM posts \n  WHERE published_at > DATE_SUB(NOW(), INTERVAL 30 DAY)\n)\nORDER BY post_count DESC"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithComments_PreservesComments(t *testing.T) {
	// Given
	input := `
	|-- This query finds active users|
	|SELECT id, name, email|
	|FROM users|
	|WHERE active = true  -- Only active users|
	|  AND created_at > '2024-01-01'  /* Users from this year */|
	|ORDER BY name|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "-- This query finds active users\nSELECT id, name, email\nFROM users\nWHERE active = true  -- Only active users\n  AND created_at > '2024-01-01'  /* Users from this year */\nORDER BY name"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithStoredProcedure_PreservesStoredProcedure(t *testing.T) {
	// Given
	input := `
	|DELIMITER //|
	|CREATE PROCEDURE GetUserStats(IN user_id INT)|
	|BEGIN|
	|  SELECT |
	|    COUNT(*) as total_posts,|
	|    MAX(created_at) as last_post_date|
	|  FROM posts |
	|  WHERE posts.user_id = user_id;|
	|END //|
	|DELIMITER ;|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "DELIMITER //\nCREATE PROCEDURE GetUserStats(IN user_id INT)\nBEGIN\n  SELECT \n    COUNT(*) as total_posts,\n    MAX(created_at) as last_post_date\n  FROM posts \n  WHERE posts.user_id = user_id;\nEND //\nDELIMITER ;"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithCTEQuery_PreservesWithClause(t *testing.T) {
	// Given
	input := `
	|WITH active_users AS (|
	|  SELECT id, name, email|
	|  FROM users|
	|  WHERE active = true|
	|),|
	|user_posts AS (|
	|  SELECT user_id, COUNT(*) as post_count|
	|  FROM posts|
	|  GROUP BY user_id|
	|)|
	|SELECT au.name, au.email, COALESCE(up.post_count, 0) as posts|
	|FROM active_users au|
	|LEFT JOIN user_posts up ON au.id = up.user_id|
	|ORDER BY posts DESC|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "WITH active_users AS (\n  SELECT id, name, email\n  FROM users\n  WHERE active = true\n),\nuser_posts AS (\n  SELECT user_id, COUNT(*) as post_count\n  FROM posts\n  GROUP BY user_id\n)\nSELECT au.name, au.email, COALESCE(up.post_count, 0) as posts\nFROM active_users au\nLEFT JOIN user_posts up ON au.id = up.user_id\nORDER BY posts DESC"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithStringLiteralsContainingPipes_PreservesLiterals(t *testing.T) {
	// Given
	input := `
	|SELECT |
	|  id,|
	|  CONCAT(first_name, '|', last_name) as full_name,|
	|  'Status: Active|Verified' as status_info|
	|FROM users|
	|WHERE description LIKE '%pipe|symbol%'|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT \n  id,\n  CONCAT(first_name, '|', last_name) as full_name,\n  'Status: Active|Verified' as status_info\nFROM users\nWHERE description LIKE '%pipe|symbol%'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithWindowFunctions_PreservesWindowSyntax(t *testing.T) {
	// Given
	input := `
	|SELECT |
	|  name,|
	|  salary,|
	|  ROW_NUMBER() OVER (ORDER BY salary DESC) as rank,|
	|  AVG(salary) OVER (PARTITION BY department) as dept_avg|
	|FROM employees|
	|ORDER BY rank|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT \n  name,\n  salary,\n  ROW_NUMBER() OVER (ORDER BY salary DESC) as rank,\n  AVG(salary) OVER (PARTITION BY department) as dept_avg\nFROM employees\nORDER BY rank"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMixedPipeContent_HandlesInternalPipesCorrectly(t *testing.T) {
	// Given
	input := `
	|SELECT id, name || ' | ' || email as display_name|
	|FROM users|
	|WHERE status IN ('active|premium', 'trial|basic')|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "SELECT id, name || ' | ' || email as display_name\nFROM users\nWHERE status IN ('active|premium', 'trial|basic')"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}
