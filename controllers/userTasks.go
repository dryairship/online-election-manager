package controllers

import (
    "github.com/dryairship/online-election-manager/db"
)

var ElectionDb db.ElectionDatabase

type UserError struct {
    reason string
}

func (err *UserError) Error() string {
    return err.reason
}

func CanMailBeSentToStudent(roll string) (bool, error) {
    voter, err := ElectionDb.FindVoter(roll)
    if err != nil {
        return true, nil
    } else {
        if voter.AuthCode == "" {
            return false, &UserError{"Student has already registered."}
        } else {
            return false, &UserError{"Verification mail has already been sent to this student."}
        }
    }
}

