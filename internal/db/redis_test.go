package db

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bliuchak/heroes/internal/storage"
	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/assert"
)

type statusExpected struct {
	isError  bool
	response string
}

func TestDbRedis_Status(t *testing.T) {
	tests := []struct {
		name      string
		redisStub radix.Client
		expected  statusExpected
	}{
		{
			name: "should return successful response PONG",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "PING":
					return "PONG"
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: statusExpected{
				response: "PONG",
			},
		},
		{
			name: "should return error message",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "PING":
					return fmt.Errorf("%q", "error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: statusExpected{
				isError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{client: tt.redisStub}
			res, err := r.Status()

			if tt.expected.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected.response, res)
		})
	}
}

type getHeroesExpected struct {
	isError bool
	error   error
	heroes  []storage.Hero
}

func TestDbRedis_GetHeroes(t *testing.T) {
	tests := []struct {
		name      string
		redisStub radix.Client
		expected  getHeroesExpected
	}{
		{
			name: "should return successfully scanned dataset response with one hero",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "SCAN":
					if cur := args[1]; cur == "0" {
						keys := []string{"hero.1"}
						return []interface{}{"1", keys}
					}
					return []interface{}{"0", []string{}}
				case "GET":
					return "Batman"
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroesExpected{
				heroes: []storage.Hero{
					{
						ID:   "1",
						Name: "Batman",
					},
				},
			},
		},
		{
			name: "should return error on GET for scanned key",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "SCAN":
					if cur := args[1]; cur == "0" {
						keys := []string{"hero.1"}
						return []interface{}{"1", keys}
					}
					return []interface{}{"0", []string{}}
				case "GET":
					return errors.New("GET error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroesExpected{
				isError: true,
				heroes:  []storage.Hero{},
				error:   errors.New("GET error"),
			},
		},
		{
			name: "should return error on scanner.Close()",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "SCAN":
					return errors.New("scanner error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroesExpected{
				isError: true,
				heroes:  []storage.Hero{},
				error:   errors.New("scanner error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{client: tt.redisStub}
			res, err := r.GetHeroes()

			if tt.expected.isError {
				assert.Equal(t, err, tt.expected.error)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected.heroes, res)
		})
	}
}

type getHeroExpected struct {
	isError bool
	error   error
	hero    storage.Hero
}

func TestDbRedis_GetHero(t *testing.T) {
	tests := []struct {
		name      string
		redisStub radix.Client
		expected  getHeroExpected
	}{
		{
			name: "should return error on EXISTS command",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "EXISTS":
					return errors.New("EXISTS error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroExpected{
				isError: true,
				error:   errors.New("EXISTS error"),
			},
		},
		{
			name: "should return error hero not existing",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "EXISTS":
					return 0
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroExpected{
				isError: true,
				error:   storage.NewErrHeroNotExist("hero not exist"),
			},
		},
		{
			name: "should return error on GET command",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "EXISTS":
					return 1
				case "GET":
					return errors.New("GET error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroExpected{
				isError: true,
				error:   errors.New("GET error"),
			},
		},
		{
			name: "should return correct hero information",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "EXISTS":
					return 1
				case "GET":
					return "Batman"
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: getHeroExpected{
				hero: storage.Hero{
					ID:   "1",
					Name: "Batman",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{client: tt.redisStub}
			res, err := r.GetHero("1")

			if tt.expected.isError {
				assert.Equal(t, err, tt.expected.error)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected.hero, res)
		})
	}
}

type createHeroExpected struct {
	isError bool
	error   error
}

func TestDbRedis_CreateHero(t *testing.T) {
	tests := []struct {
		name      string
		redisStub radix.Client
		expected  createHeroExpected
	}{
		{
			name: "should return error on SET command",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "SET":
					return errors.New("SET error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: createHeroExpected{
				isError: true,
				error:   errors.New("SET error"),
			},
		},
		{
			name: "should return no errors",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "SET":
					return nil
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: createHeroExpected{
				isError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{client: tt.redisStub}
			err := r.CreateHero("1", "Batman")

			if tt.expected.isError {
				assert.Equal(t, err, tt.expected.error)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type deleteHeroExpected struct {
	isError bool
	error   error
}

func TestDbRedis_DeleteHero(t *testing.T) {
	tests := []struct {
		name      string
		redisStub radix.Client
		expected  deleteHeroExpected
	}{
		{
			name: "should return error on DEL command",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "DEL":
					return errors.New("DEL error")
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: deleteHeroExpected{
				isError: true,
				error:   errors.New("DEL error"),
			},
		},
		{
			name: "should return no errors",
			redisStub: radix.Stub("", "", func(args []string) interface{} {
				switch args[0] {
				case "DEL":
					return nil
				default:
					return fmt.Errorf("testStub doesn't support command %q", args[0])
				}
			}),
			expected: deleteHeroExpected{
				isError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Redis{client: tt.redisStub}
			err := r.DeleteHero("1")

			if tt.expected.isError {
				assert.Equal(t, err, tt.expected.error)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDb_NewPool(t *testing.T) {
	_, err := NewRedis("", "", "0")
	assert.Error(t, err)
}
