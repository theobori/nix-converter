package json

import (
	"testing"
)

var toNixData = []string{
	`{
  "globalOperations": {
    "regions": {
      "northAmerica": {
        "countries": {
          "unitedStates": {
            "states": {
              "california": {
                "cities": {
                  "sanFrancisco": {
                    "demographics": {
                      "population": 874961,
                      "ageDistribution": {
                        "under18": 13.4,
                        "18to24": 8.9,
                        "25to44": 45.2,
                        "45to64": 20.7,
                        "over65": 11.8
                      },
                      "householdIncome": {
                        "median": 112449,
                        "distribution": {
                          "under50k": 23.4,
                          "50kto100k": 25.1,
                          "100kto150k": 18.7,
                          "150kto200k": 11.2,
                          "over200k": 21.6
                        }
                      }
                    },
                    "infrastructure": {
                      "transportation": {
                        "public": {
                          "subway": {
                            "lines": ["red", "blue", "green"],
                            "stations": 45,
                            "dailyRidership": 157980,
                            "maintenance": {
                              "lastInspection": "2024-11-15",
                              "nextScheduled": "2025-01-15",
                              "issues": {
                                "critical": 0,
                                "moderate": 3,
                                "minor": 12
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          },
          "canada": {
            "provinces": {
              "ontario": {
                "cities": {
                  "toronto": {
                    "businesses": {
                      "technology": {
                        "startups": {
                          "company1": {
                            "details": {
                              "founded": 2022,
                              "employees": 45,
                              "funding": {
                                "rounds": {
                                  "seed": {
                                    "amount": 1500000,
                                    "investors": {
                                      "lead": {
                                        "name": "Tech Ventures",
                                        "stake": 15,
                                        "board": {
                                          "seats": 2,
                                          "members": [
                                            {
                                              "name": "Jane Smith",
                                              "position": "Managing Partner",
                                              "experience": {
                                                "years": 15,
                                                "expertise": ["AI", "SaaS", "Fintech"],
                                                "previous": {
                                                  "companies": [
                                                    {
                                                      "name": "Growth Capital",
                                                      "role": "Partner",
                                                      "duration": 8,
                                                      "investments": {
                                                        "successful": 12,
                                                        "total": 15
                                                      }
                                                    }
                                                  ]
                                                }
                                              }
                                            }
                                          ]
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      },
      "europe": {
        "countries": {
          "germany": {
            "states": {
              "bavaria": {
                "industries": {
                  "automotive": {
                    "manufacturers": {
                      "company1": {
                        "production": {
                          "facilities": {
                            "mainPlant": {
                              "capacity": {
                                "daily": 1200,
                                "models": {
                                  "sedan": {
                                    "variants": {
                                      "standard": {
                                        "specifications": {
                                          "engine": {
                                            "type": "hybrid",
                                            "power": 245,
                                            "efficiency": {
                                              "city": 52,
                                              "highway": 48,
                                              "combined": 50
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`,
	`{"hello": "world"}`,
}

func TestJSONToNix(t *testing.T) {
	for _, s := range toNixData {
		nj := NewNixJson(s)

		_, err := nj.ToNix()
		if err != nil {
			t.Fatal(err)
		}
	}
}
