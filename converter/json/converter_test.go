package json

import (
	"testing"
	// "github.com/theobori/nix-converter/internal/common"
)

var jsonStrings = []string{
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
                            "lines": [
                              "red",
                              "blue",
                              "green"
                            ],
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
                                                "expertise": [
                                                  "AI",
                                                  "SaaS",
                                                  "Fintech"
                                                ],
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
                                              "combined": 50,
                                              "": 123,
                                              "123": ""
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
}

var nixStrings = []string{
	`{
  id = "c7d8e9f0";
  users = [
    {
      name = "Alice";
      age = 28;
      pets = [
        {
          type = "cat";
          name = "Luna";
          toys = [
            "mouse"
            "ball"
          ];
        }
        {
          type = "dog";
          name = "Max";
        }
      ];
    }
    {
      name = "Bob";
      age = 34;
      age2 = -34;
      age3 = -3.45;
      ag2 = -0;
      "" = 123;
      hello = "";
      age12 = -0.45;
      ag2123 = 0.001;
      age4 = -0.0045;
      age5 = -0.00000000000000000000000045;
      pets = null;
    }
  ];
  settings = {
    theme = {
      dark = {
        primary = "#1a1a1a";
        accent = "#4287f5";
      };
      light = {
        primary = "#ffffff";
        accent = "#2196f3";
      };
    };
    notifications = true;
  };
  meta = {
    created = "2024-01-01";
    modified = {
      by = "system";
      timestamp = "2024-02-15T14:30:00Z";
    };
  };
}`,
}

func TestJSONToNix(t *testing.T) {
	// common.TestHelperToNixStrings(t, jsonStrings, FromNix, ToNix)
}

func TestJSONFromNix(t *testing.T) {
	// common.TestHelperFromNixStrings(t, nixStrings, FromNix, ToNix)
}
