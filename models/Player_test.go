package models

import (
	"strings"
	"testing"
)

func TestPlayer_Validate(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantErr   bool
		errMsg    string
	}{
		// Valid cases
		{
			name:      "valid names",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   false,
		},
		{
			name:      "minimum length names",
			firstName: "Al",
			lastName:  "Li",
			wantErr:   false,
		},
		{
			name:      "maximum length names",
			firstName: strings.Repeat("a", 20),
			lastName:  strings.Repeat("b", 20),
			wantErr:   false,
		},
		{
			name:      "names with mixed case",
			firstName: "JoHn",
			lastName:  "DoE",
			wantErr:   false,
		},

		// First name too short
		{
			name:      "first name too short - empty",
			firstName: "",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},
		{
			name:      "first name too short - one character",
			firstName: "J",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},

		// First name too long
		{
			name:      "first name too long - 50 characters",
			firstName: strings.Repeat("a", 50),
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},

		// Last name too short
		{
			name:      "last name too short - empty",
			firstName: "John",
			lastName:  "",
			wantErr:   true,
			errMsg:    "Last name must be between 2 and 20",
		},
		{
			name:      "last name too short - one character",
			firstName: "John",
			lastName:  "D",
			wantErr:   true,
			errMsg:    "Last name must be between 2 and 20",
		},

		// Last name too long
		{
			name:      "last name too long - 50 characters",
			firstName: "John",
			lastName:  strings.Repeat("b", 50),
			wantErr:   true,
			errMsg:    "Last name must be between 2 and 20",
		},

		// Multiple parts in first name
		{
			name:      "first name with space",
			firstName: "John Paul",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with multiple spaces",
			firstName: "John Paul George",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with leading space",
			firstName: " John",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with trailing space",
			firstName: "John ",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with tabs",
			firstName: "John\tPaul",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},

		// Multiple parts in last name
		{
			name:      "last name with space",
			firstName: "John",
			lastName:  "Van Dijk",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "last name with multiple spaces",
			firstName: "John",
			lastName:  "De La Cruz",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "last name with leading space",
			firstName: "John",
			lastName:  " Doe",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "last name with trailing space",
			firstName: "John",
			lastName:  "Doe ",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "last name with newline",
			firstName: "John",
			lastName:  "Doe\nSmith",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},

		// Both names invalid
		{
			name:      "both names too short",
			firstName: "J",
			lastName:  "D",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},
		{
			name:      "both names with spaces",
			firstName: "John Paul",
			lastName:  "Van Dijk",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "both names empty",
			firstName: "",
			lastName:  "",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},

		// Edge cases with whitespace only
		{
			name:      "first name only spaces",
			firstName: "   ",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between 2 and 20",
		},
		{
			name:      "last name only spaces",
			firstName: "John",
			lastName:  "   ",
			wantErr:   true,
			errMsg:    "Last name must be between 2 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}

			err := p.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPlayer_SetPlayerName(t *testing.T) {
	tests := []struct {
		name           string
		firstName      string
		lastName       string
		wantPlayerName string
		wantFirstName  string
		wantLastName   string
	}{
		{
			name:           "normal names",
			firstName:      "John",
			lastName:       "Doe",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with leading spaces",
			firstName:      "  John",
			lastName:       "  Doe",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with trailing spaces",
			firstName:      "John  ",
			lastName:       "Doe  ",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with both leading and trailing spaces",
			firstName:      "  John  ",
			lastName:       "  Doe  ",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with tabs",
			firstName:      "\tJohn\t",
			lastName:       "\tDoe\t",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with newlines",
			firstName:      "\nJohn\n",
			lastName:       "\nDoe\n",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
		{
			name:           "names with mixed whitespace",
			firstName:      " \t\nJohn \t\n",
			lastName:       " \t\nDoe \t\n",
			wantPlayerName: "John Doe",
			wantFirstName:  "John",
			wantLastName:   "Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{}
			p.SetPlayerName(tt.firstName, tt.lastName)

			if p.PlayerName != tt.wantPlayerName {
				t.Errorf("SetPlayerName() PlayerName = %q, want %q", p.PlayerName, tt.wantPlayerName)
			}
			if p.FirstName != tt.wantFirstName {
				t.Errorf("SetPlayerName() FirstName = %q, want %q", p.FirstName, tt.wantFirstName)
			}
			if p.LastName != tt.wantLastName {
				t.Errorf("SetPlayerName() LastName = %q, want %q", p.LastName, tt.wantLastName)
			}
		})
	}
}

func TestPlayer_ValidateFirstAndLastNameLength(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid lengths",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   false,
		},
		{
			name:      "minimum valid lengths",
			firstName: "Jo",
			lastName:  "Do",
			wantErr:   false,
		},
		{
			name:      "maximum valid lengths",
			firstName: strings.Repeat("a", 20),
			lastName:  strings.Repeat("b", 20),
			wantErr:   false,
		},
		{
			name:      "first name below minimum",
			firstName: "J",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between",
		},
		{
			name:      "first name above maximum",
			firstName: strings.Repeat("a", 21),
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must be between",
		},
		{
			name:      "last name below minimum",
			firstName: "John",
			lastName:  "D",
			wantErr:   true,
			errMsg:    "Last name must be between",
		},
		{
			name:      "last name above maximum",
			firstName: "John",
			lastName:  strings.Repeat("b", 21),
			wantErr:   true,
			errMsg:    "Last name must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}

			err := p.validateFirstAndLastNameLength()

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateFirstAndLastNameLength() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFirstAndLastNameLength() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateFirstAndLastNameLength() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPlayer_ValidateFirstAndLastNameHasOnePart(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid single part names",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   false,
		},
		{
			name:      "first name with space",
			firstName: "John Paul",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with multiple spaces",
			firstName: "John Paul George",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "last name with space",
			firstName: "John",
			lastName:  "Van Dijk",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "last name with multiple spaces",
			firstName: "John",
			lastName:  "De La Cruz",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
		{
			name:      "both names with spaces",
			firstName: "John Paul",
			lastName:  "Van Dijk",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "first name with tab",
			firstName: "John\tPaul",
			lastName:  "Doe",
			wantErr:   true,
			errMsg:    "First name must have exactly 1 part",
		},
		{
			name:      "last name with newline",
			firstName: "John",
			lastName:  "Doe\nSmith",
			wantErr:   true,
			errMsg:    "Last name must have exactly 1 part",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}

			err := p.validateFirstAndLastNameHasOnePart()

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateFirstAndLastNameHasOnePart() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFirstAndLastNameHasOnePart() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateFirstAndLastNameHasOnePart() unexpected error = %v", err)
				}
			}
		})
	}
}
