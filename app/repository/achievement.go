package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"uas/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	CreateAchievementMongo(ctx context.Context, data models.AchievementMongo) (string, error)
	CreateAchievementReference(ctx context.Context, ref models.AchievementReference) error
	GetAchievementByID(ctx context.Context, id string) (models.AchievementReference, error)
    UpdateAchievement(ctx context.Context, pgID string, mongoID string, data models.AchievementMongo) error
    SoftDeleteAchievement(ctx context.Context, pgID string, mongoID string) error
	SubmitAchievement(ctx context.Context, id string) error
    GetLecturerIDByUserID(ctx context.Context, userID string) (string, error)
    VerifyAchievement(ctx context.Context, id string, verifierUserID string) error
    RejectAchievement(ctx context.Context, id string, verifierUserID string, note string) error
    CheckStudentAdvisorRelationship(ctx context.Context, lecturerID string, studentID string) (bool, error)
    GetAllReferences(ctx context.Context) ([]models.AchievementReference, map[string]string, map[string]string, error)
    GetMongoDetailsByIDs(ctx context.Context, mongoIDs []string) (map[string]models.AchievementMongo, error)
    GetAchievementReferenceWithDetail(ctx context.Context, id string) (models.AchievementResponse, error)
    GetMongoDetailByID(ctx context.Context, mongoID string) (models.AchievementMongo, error)
}

type achievementRepository struct {
	pg    *sql.DB
	mongo *mongo.Database
}

func NewAchievementRepository(pg *sql.DB, mongo *mongo.Database) AchievementRepository {
	return &achievementRepository{pg: pg, mongo: mongo}
}

func (r *achievementRepository) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT id FROM students WHERE user_id = $1`
	var studentID string
	err := r.pg.QueryRowContext(ctx, query, userID).Scan(&studentID)
	if err != nil {
		return "", err
	}
	return studentID, nil
}

// Simpan ke MongoDB
func (r *achievementRepository) CreateAchievementMongo(ctx context.Context, data models.AchievementMongo) (string, error) {
	collection := r.mongo.Collection("achievements")
	
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return "", fmt.Errorf("gagal insert ke mongo: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("gagal cast insertedID")
	}

	return oid.Hex(), nil
}

// Simpan Referensi (PostgreSQL)
func (r *achievementRepository) CreateAchievementReference(ctx context.Context, ref models.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (
			id, student_id, mongo_achievement_id, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $5)
	`
	_, err := r.pg.ExecContext(ctx, query, ref.ID, ref.StudentID, ref.MongoAchievementID, "draft", time.Now())
	if err != nil {
		return fmt.Errorf("gagal insert ke postgres: %w", err)
	}
	return nil
}

// Ambil Data Achievement berdasarkan ID (Postgres)
func (r *achievementRepository) GetAchievementByID(ctx context.Context, id string) (models.AchievementReference, error) {
    query := `
        SELECT id, student_id, mongo_achievement_id, status 
        FROM achievement_references 
        WHERE id = $1 AND deleted_at IS NULL
    `
    var ref models.AchievementReference    
    err := r.pg.QueryRowContext(ctx, query, id).Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status)
    if err != nil {
        return models.AchievementReference{}, err
    }
    return ref, nil
}

// Update Achievement (Mongo & Postgres Timestamp)
func (r *achievementRepository) UpdateAchievement(ctx context.Context, pgID string, mongoID string, data models.AchievementMongo) error {
    // Update MongoDB (Detail)
    mongoOID, _ := primitive.ObjectIDFromHex(mongoID)
    filter := bson.M{"_id": mongoOID}
    
    update := bson.M{
        "$set": bson.M{
            "achievementType": data.AchievementType,
            "title":           data.Title,
            "description":     data.Description,
            "details":         data.Details,
            "tags":            data.Tags,
            "updatedAt":       time.Now(),
        },
    }

    _, err := r.mongo.Collection("achievements").UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("gagal update mongo: %w", err)
    }

    // Update PostgreSQL (Hanya updated_at)
    queryPG := `UPDATE achievement_references SET updated_at = NOW() WHERE id = $1`
    _, err = r.pg.ExecContext(ctx, queryPG, pgID)
    if err != nil {
        return fmt.Errorf("gagal update postgres: %w", err)
    }

    return nil
}

// Soft Delete
func (r *achievementRepository) SoftDeleteAchievement(ctx context.Context, pgID string, mongoID string) error {
    // Soft Delete Postgres
    queryPG := `UPDATE achievement_references SET deleted_at = NOW() WHERE id = $1`
    _, err := r.pg.ExecContext(ctx, queryPG, pgID)
    if err != nil {
        return fmt.Errorf("gagal soft delete postgres: %w", err)
    }

    // Soft Delete Mongo
    mongoOID, _ := primitive.ObjectIDFromHex(mongoID)
    filter := bson.M{"_id": mongoOID}
    update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
    
    _, err = r.mongo.Collection("achievements").UpdateOne(ctx, filter, update)
    return err
}

func (r *achievementRepository) SubmitAchievement(ctx context.Context, id string) error {
    query := `
        UPDATE achievement_references 
        SET status = 'submitted', 
            submitted_at = NOW(), 
            updated_at = NOW() 
        WHERE id = $1
    `

    result, err := r.pg.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("gagal submit prestasi: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("data tidak ditemukan atau tidak ada perubahan")
    }

    return nil
}

func (r *achievementRepository) GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
    query := `SELECT id FROM lecturers WHERE user_id = $1`
    var lecturerID string
    err := r.pg.QueryRowContext(ctx, query, userID).Scan(&lecturerID)
    if err != nil {
        return "", err
    }
    return lecturerID, nil
}

func (r *achievementRepository) VerifyAchievement(ctx context.Context, id string, verifierUserID string) error {
    query := `
        UPDATE achievement_references 
        SET status = 'verified', 
            verified_by = $2, 
            verified_at = NOW(),
            updated_at = NOW()
        WHERE id = $1
    `
    result, err := r.pg.ExecContext(ctx, query, id, verifierUserID)
    if err != nil {
        return fmt.Errorf("gagal verifikasi: %w", err)
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("data tidak ditemukan")
    }
    return nil
}

func (r *achievementRepository) RejectAchievement(ctx context.Context, id string, verifierUserID string, note string) error {
    query := `
        UPDATE achievement_references 
        SET status = 'rejected', 
            verified_by = $2, 
            rejection_note = $3,
            updated_at = NOW()
        WHERE id = $1
    `
    result, err := r.pg.ExecContext(ctx, query, id, verifierUserID, note)
    if err != nil {
        return fmt.Errorf("gagal reject: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("data tidak ditemukan")
    }
    return nil
}

func (r *achievementRepository) CheckStudentAdvisorRelationship(ctx context.Context, lecturerID string, studentID string) (bool, error) {
    query := `SELECT count(1) FROM students WHERE id = $1 AND advisor_id = $2`
    
    var count int
    err := r.pg.QueryRowContext(ctx, query, studentID, lecturerID).Scan(&count)
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}

func (r *achievementRepository) GetAllReferences(ctx context.Context) ([]models.AchievementReference, map[string]string, map[string]string, error) {
    query := `
        SELECT 
            ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.created_at,
            u.full_name, s.student_id as nim
        FROM achievement_references ar
        JOIN students s ON ar.student_id = s.id
        JOIN users u ON s.user_id = u.id
        ORDER BY ar.created_at DESC
    `
    
    rows, err := r.pg.QueryContext(ctx, query)
    if err != nil {
        return nil, nil, nil, err
    }
    defer rows.Close()

    var refs []models.AchievementReference
    studentNames := make(map[string]string)
    studentNIMs := make(map[string]string)

    for rows.Next() {
        var ref models.AchievementReference
        var fullName, nim string
        
        err := rows.Scan(
            &ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &ref.CreatedAt,
            &fullName, &nim,
        )
        if err != nil {
            return nil, nil, nil, err
        }

        refs = append(refs, ref)
        studentNames[ref.ID] = fullName
        studentNIMs[ref.ID] = nim
    }

    return refs, studentNames, studentNIMs, nil
}

func (r *achievementRepository) GetMongoDetailsByIDs(ctx context.Context, mongoIDs []string) (map[string]models.AchievementMongo, error) {
    var objectIDs []primitive.ObjectID
    for _, id := range mongoIDs {
        oid, err := primitive.ObjectIDFromHex(id)
        if err == nil {
            objectIDs = append(objectIDs, oid)
        }
    }

    filter := bson.M{"_id": bson.M{"$in": objectIDs}}
    cursor, err := r.mongo.Collection("achievements").Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    results := make(map[string]models.AchievementMongo)
    for cursor.Next(ctx) {
        var doc models.AchievementMongo
        if err := cursor.Decode(&doc); err == nil {
            results[doc.ID.Hex()] = doc
        }
    }
    return results, nil
}

func (r *achievementRepository) GetAchievementReferenceWithDetail(ctx context.Context, id string) (models.AchievementResponse, error) {
    query := `
        SELECT 
            ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, 
            ar.created_at, ar.submitted_at, ar.verified_at, ar.rejection_note,
            u.full_name, s.student_id as nim
        FROM achievement_references ar
        JOIN students s ON ar.student_id = s.id
        JOIN users u ON s.user_id = u.id
        WHERE ar.id = $1 AND ar.deleted_at IS NULL
    `
    
    var res models.AchievementResponse
    var submittedAt, verifiedAt *time.Time
    var rejectionNote *string

    err := r.pg.QueryRowContext(ctx, query, id).Scan(
        &res.ID, &res.StudentID, &res.MongoID, &res.Status,
        &res.CreatedAt, &submittedAt, &verifiedAt, &rejectionNote,
        &res.StudentName, &res.StudentNIM,
    )
    if err != nil {
        return models.AchievementResponse{}, err
    }

    res.SubmittedAt = submittedAt
    res.VerifiedAt = verifiedAt
    if rejectionNote != nil {
        res.RejectionNote = *rejectionNote
    }

    return res, nil
}

func (r *achievementRepository) GetMongoDetailByID(ctx context.Context, mongoID string) (models.AchievementMongo, error) {
    oid, err := primitive.ObjectIDFromHex(mongoID)
    if err != nil {
        return models.AchievementMongo{}, err
    }

    var result models.AchievementMongo
    err = r.mongo.Collection("achievements").FindOne(ctx, bson.M{"_id": oid}).Decode(&result)
    if err != nil {
        return models.AchievementMongo{}, err
    }

    return result, nil
}