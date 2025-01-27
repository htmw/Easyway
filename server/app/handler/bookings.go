package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/ksharma67/EasyWay/server/app/model"
)

//Booking API

func CreateBooking(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println("Logged CreateBooking:POST")
	booking := model.Booking{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&booking); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	var bookings []model.Booking
	db.Where("service_id = ? ", booking.ServiceId).Find(&bookings)

	for _, b := range bookings {
		start := b.StartTime
		end := b.EndTime
		date := b.Date
		fmt.Println("date=" + date + " start=" + start + " end=" + end)
		fmt.Println("date=" + booking.Date + " start=" + booking.StartTime + " end=" + booking.EndTime)
		if booking.Date == date && booking.EndTime <= end && booking.StartTime >= start {
			respondError(w, http.StatusInternalServerError, "Time slot unavailable")
			return
		}
	}

	if err := db.Save(&booking).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, booking)
}

//Booking Cancellation API
func GetCancelledBookings(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()["userId"]
	i, err := strconv.Atoi(userId[0])
	if err == nil {
		fmt.Print("Good")
	}
	cancelledBookings := []model.Booking{}
	if err := db.Where("user_id = ? AND is_cancelled = ?", i, true).Find(&cancelledBookings).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	}

	respondJSON(w, http.StatusOK, cancelledBookings)

}

func CancelBooking(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	bookingId := r.URL.Query()["id"]
	i, err := strconv.Atoi(bookingId[0])
	if err != nil {
		fmt.Println(err)
	}

	booking := model.Booking{}

	// err := r.ParseForm()
	// if err != nil {
	// 	respondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	// b := model.Booking{}

	// decoder := json.NewDecoder(r.Body)

	// if err := decoder.Decode(&b); err != nil {
	// 	respondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
	// defer r.Body.Close()

	db.Where("id = ?", i).Find(&booking)

	if booking.IsCancelled == true {
		respondJSON(w, http.StatusAlreadyReported, "Booking already cancelled")
		return
	}

	db.Model(&booking).Where("id = ?", i).Update("is_cancelled", true)
	respondJSON(w, http.StatusOK, "Booking is cancelled")

}

func GetBookings(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()["userId"]
	i, err := strconv.Atoi(userId[0])
	if err != nil {
		fmt.Println(err)
	}

	booking := []model.Booking{}
	if err := db.Where("user_id = ? AND is_cancelled = ?", i, false).Find(&booking).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	}

	respondJSON(w, http.StatusOK, booking)
}

// getUserOr404 gets a booking instance if exists, or respond the 404 error otherwise
func getBookingInfo(db *gorm.DB, custId int, w http.ResponseWriter, r *http.Request) *model.Booking {
	booking := model.Booking{}
	return &booking
}


// Edit Booking API
func UpdateBooking(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    fmt.Println("Logged UpdateBooking:PUT")
    booking := model.Booking{}

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&booking); err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }
    defer r.Body.Close()

    // Check if the booking exists in the database
    existingBooking := model.Booking{}
    if err := db.First(&existingBooking, booking.Id).Error; err != nil {
        respondError(w, http.StatusNotFound, "Booking not found")
        return
    }

    // Check if the new booking time slot is available
    var conflictingBookings []model.Booking
    db.Where("service_id = ? AND id != ? AND date = ?", booking.ServiceId, booking.Id, booking.Date).
        Find(&conflictingBookings)
    for _, b := range conflictingBookings {
        start := b.StartTime
        end := b.EndTime
        date := b.Date
        fmt.Println("date=" + date + " start=" + start + " end=" + end)
        fmt.Println("date=" + booking.Date + " start=" + booking.StartTime + " end=" + booking.EndTime)
        if booking.Date == date && booking.EndTime <= end && booking.StartTime >= start {
            respondError(w, http.StatusInternalServerError, "Time slot unavailable")
            return
        }
    }

    // Update the booking record in the database
    existingBooking.ServiceId = booking.ServiceId
    existingBooking.Date = booking.Date
    existingBooking.StartTime = booking.StartTime
    existingBooking.EndTime = booking.EndTime
    existingBooking.Note = booking.Note
    if err := db.Save(&existingBooking).Error; err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, existingBooking)
}
