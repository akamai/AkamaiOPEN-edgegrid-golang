package ivm

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPolicies(t *testing.T) {
	tests := map[string]struct {
		params           ListPoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListPoliciesResponse
		expectedHeaders  map[string][]string
		withError        error
	}{
		"200 OK": {
			params: ListPoliciesRequest{
				Network:     PolicyNetworkStaging,
				Contract:    "3-WNKXX1",
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "itemKind": "POLICY",
    "items": [
        {
            "id": ".auto",
            "version": 1,
            "previousVersion": 0,
            "rolloutInfo": {
                "startTime": 1626379176,
                "endTime": 1626379177,
                "rolloutDuration": 1
            },
            "video": false,
            "user": "system",
            "dateCreated": "2021-07-15 19:59:35+0000"
        },
        {
            "id": "asd",
            "version": 2,
            "previousVersion": 1,
            "rolloutInfo": {
                "startTime": 1638894035,
                "endTime": 1638894036,
                "rolloutDuration": 1
            },
            "breakpoints": {
                "widths": [
                    320,
                    640,
                    1024,
                    2048,
                    5000
                ]
            },
            "output": {
                "perceptualQuality": "mediumHigh"
            },
            "transformations": [
                {
                    "transformation": "Composite",
                    "xPosition": 0,
                    "yPosition": 0,
                    "gravity": "NorthWest",
                    "placement": "Over",
                    "image": {
                        "type": "Text",
                        "fill": "#000000",
                        "size": 72,
                        "stroke": "#FFFFFF",
                        "strokeSize": 0,
                        "text": "Hello There",
                        "transformation": {
                            "transformation": "Compound",
                            "transformations": []
                        }
                    }
                }
            ],
            "video": false,
            "user": "jsmith",
            "dateCreated": "2021-12-07 16:20:34+0000"
        },
        {
            "id": "cheese",
            "version": 1,
            "previousVersion": 0,
            "rolloutInfo": {
                "startTime": 1628279193,
                "endTime": 1628279194,
                "rolloutDuration": 1
            },
            "breakpoints": {
                "widths": [
                    720,
                    1080,
                    1366,
                    1920,
                    5000
                ]
            },
            "output": {
                "perceptualQuality": "mediumHigh"
            },
            "video": false,
            "user": "jsmith",
            "dateCreated": "2021-08-06 19:46:32+0000"
        },
        {
            "id": "example",
            "version": 9,
            "previousVersion": 8,
            "rolloutInfo": {
                "startTime": 1639680399,
                "endTime": 1639680400,
                "rolloutDuration": 1
            },
            "breakpoints": {
                "widths": [
                    320,
                    640,
                    1024,
                    2048,
                    5000
                ]
            },
            "output": {
                "perceptualQuality": "mediumHigh"
            },
            "transformations": [
                {
                    "transformation": "Blur",
                    "sigma": {
                        "var": "blurVar"
                    }
                },
                {
                    "transformation": "MaxColors",
                    "colors": 2
                }
            ],
            "variables": [
                {
                    "name": "blurVar",
                    "type": "number",
                    "defaultValue": "5"
                }
            ],
            "video": false,
            "user": "foofoo5",
            "dateCreated": "2021-12-16 18:46:38+0000"
        },
        {
            "id": "mgw",
            "version": 1,
            "previousVersion": 0,
            "rolloutInfo": {
                "startTime": 1639680457,
                "endTime": 1639680458,
                "rolloutDuration": 1
            },
            "breakpoints": {
                "widths": [
                    320,
                    640,
                    1024,
                    2048,
                    5000
                ]
            },
            "output": {
                "perceptualQuality": "mediumHigh"
            },
            "transformations": [
                {
                    "transformation": "Blur",
                    "sigma": {
                        "var": "blurVar"
                    }
                },
                {
                    "transformation": "MaxColors",
                    "colors": 2
                }
            ],
            "variables": [
                {
                    "name": "blurVar",
                    "type": "number",
                    "defaultValue": "5"
                }
            ],
            "video": false,
            "user": "foofoo5",
            "dateCreated": "2021-12-16 18:47:36+0000"
        },
        {
            "id": "testPolicy2",
            "version": 1,
            "previousVersion": 0,
            "rolloutInfo": {
                "startTime": 1643052400,
                "endTime": 1643052401,
                "rolloutDuration": 1
            },
            "video": false,
            "user": "jsmith",
            "dateCreated": "2022-01-24 19:26:39+0000"
        },
        {
            "id": "testpolicy",
            "version": 1,
            "previousVersion": 0,
            "rolloutInfo": {
                "startTime": 1643052164,
                "endTime": 1643052165,
                "rolloutDuration": 1
            },
            "video": false,
            "user": "jsmith",
            "dateCreated": "2022-01-24 19:22:43+0000"
        },
        {
            "id": "updatePolicy",
            "version": 2,
            "previousVersion": 1,
            "rolloutInfo": {
                "startTime": 1643055431,
                "endTime": 1643055432,
                "rolloutDuration": 1
            },
            "output": {
                "perceptualQuality": "mediumHigh"
            },
            "video": true,
            "user": "jsmith",
            "dateCreated": "2022-01-24 20:17:10+0000"
        }
    ],
    "totalItems": 8
}`,
			expectedPath: "/imaging/v2/network/staging/policies/",
			expectedHeaders: map[string][]string{
				"Contract":   {"3-WNKXX1"},
				"Policy-Set": {"570f9090-5dbe-11ec-8a0a-71665789c1d8"},
			},
			expectedResponse: &ListPoliciesResponse{
				ItemKind: "POLICY",
				Items: PolicyOutputs{
					&PolicyOutputImage{
						ID:              ".auto",
						Version:         1,
						PreviousVersion: 0,
						RolloutInfo: &RolloutInfo{
							StartTime:       1626379176,
							EndTime:         1626379177,
							RolloutDuration: 1,
						},
						Video:       false,
						User:        "system",
						DateCreated: "2021-07-15 19:59:35+0000",
					},
					&PolicyOutputImage{
						ID:              "asd",
						Version:         2,
						PreviousVersion: 1,
						RolloutInfo: &RolloutInfo{
							StartTime:       1638894035,
							EndTime:         1638894036,
							RolloutDuration: 1,
						},
						Breakpoints: &Breakpoints{
							Widths: []int{320, 640, 1024, 2048, 5000},
						},
						Output: &OutputImage{
							PerceptualQuality: &OutputImagePerceptualQualityVariableInline{
								Value: OutputImagePerceptualQualityPtr(OutputImagePerceptualQualityMediumHigh),
							},
						},
						Transformations: []TransformationType{
							&Composite{
								Transformation: "Composite",
								XPosition: &IntegerVariableInline{
									Value: tools.IntPtr(0),
								},
								YPosition: &IntegerVariableInline{
									Value: tools.IntPtr(0),
								},
								Gravity: &GravityVariableInline{
									Value: GravityPtr(GravityNorthWest),
								},
								Placement: &CompositePlacementVariableInline{
									Value: CompositePlacementPtr(CompositePlacementOver),
								},
								Image: &TextImageType{
									Type: "Text",
									Fill: &StringVariableInline{
										Value: tools.StringPtr("#000000"),
									},
									Size: &NumberVariableInline{
										Value: tools.Float32Ptr(72),
									},
									Stroke: &StringVariableInline{
										Value: tools.StringPtr("#FFFFFF"),
									},
									StrokeSize: &NumberVariableInline{
										Value: tools.Float32Ptr(0),
									},
									Text: &StringVariableInline{
										Value: tools.StringPtr("Hello There"),
									},
									Transformation: &Compound{
										Transformation: "Compound",
									},
								},
							},
						},
						Video:       false,
						User:        "jsmith",
						DateCreated: "2021-12-07 16:20:34+0000",
					},
					&PolicyOutputImage{
						ID:              "cheese",
						Version:         1,
						PreviousVersion: 0,
						RolloutInfo: &RolloutInfo{
							StartTime:       1628279193,
							EndTime:         1628279194,
							RolloutDuration: 1,
						},
						Breakpoints: &Breakpoints{
							Widths: []int{720, 1080, 1366, 1920, 5000},
						},
						Output: &OutputImage{
							PerceptualQuality: &OutputImagePerceptualQualityVariableInline{
								Value: OutputImagePerceptualQualityPtr(OutputImagePerceptualQualityMediumHigh),
							},
						},
						Video:       false,
						User:        "jsmith",
						DateCreated: "2021-08-06 19:46:32+0000",
					},
					&PolicyOutputImage{
						ID:              "example",
						Version:         9,
						PreviousVersion: 8,
						RolloutInfo: &RolloutInfo{
							StartTime:       1639680399,
							EndTime:         1639680400,
							RolloutDuration: 1,
						},
						Breakpoints: &Breakpoints{
							Widths: []int{320, 640, 1024, 2048, 5000},
						},
						Output: &OutputImage{
							PerceptualQuality: &OutputImagePerceptualQualityVariableInline{
								Value: OutputImagePerceptualQualityPtr(OutputImagePerceptualQualityMediumHigh),
							},
						},
						Transformations: []TransformationType{
							&Blur{
								Transformation: "Blur",
								Sigma: &NumberVariableInline{
									Name: tools.StringPtr("blurVar"),
								},
							},
							&MaxColors{
								Transformation: "MaxColors",
								Colors: &IntegerVariableInline{
									Value: tools.IntPtr(2),
								},
							},
						},
						Variables: []Variable{
							{
								Name:         "blurVar",
								Type:         VariableTypeNumber,
								DefaultValue: "5",
							},
						},
						Video:       false,
						User:        "foofoo5",
						DateCreated: "2021-12-16 18:46:38+0000",
					},
					&PolicyOutputImage{
						ID:              "mgw",
						Version:         1,
						PreviousVersion: 0,
						RolloutInfo: &RolloutInfo{
							StartTime:       1639680457,
							EndTime:         1639680458,
							RolloutDuration: 1,
						},
						Breakpoints: &Breakpoints{
							Widths: []int{320, 640, 1024, 2048, 5000},
						},
						Output: &OutputImage{
							PerceptualQuality: &OutputImagePerceptualQualityVariableInline{
								Value: OutputImagePerceptualQualityPtr(OutputImagePerceptualQualityMediumHigh),
							},
						},
						Transformations: []TransformationType{
							&Blur{
								Transformation: "Blur",
								Sigma: &NumberVariableInline{
									Name: tools.StringPtr("blurVar"),
								},
							},
							&MaxColors{
								Transformation: "MaxColors",
								Colors: &IntegerVariableInline{
									Value: tools.IntPtr(2),
								},
							},
						},
						Variables: []Variable{
							{
								Name:         "blurVar",
								Type:         VariableTypeNumber,
								DefaultValue: "5",
							},
						},
						Video:       false,
						User:        "foofoo5",
						DateCreated: "2021-12-16 18:47:36+0000",
					},
					&PolicyOutputImage{
						ID:              "testPolicy2",
						Version:         1,
						PreviousVersion: 0,
						RolloutInfo: &RolloutInfo{
							StartTime:       1643052400,
							EndTime:         1643052401,
							RolloutDuration: 1,
						},
						Video:       false,
						User:        "jsmith",
						DateCreated: "2022-01-24 19:26:39+0000",
					},
					&PolicyOutputImage{
						ID:              "testpolicy",
						Version:         1,
						PreviousVersion: 0,
						RolloutInfo: &RolloutInfo{
							StartTime:       1643052164,
							EndTime:         1643052165,
							RolloutDuration: 1,
						},
						Video:       false,
						User:        "jsmith",
						DateCreated: "2022-01-24 19:22:43+0000",
					},
					&PolicyOutputVideo{
						ID:              "updatePolicy",
						Version:         2,
						PreviousVersion: 1,
						RolloutInfo: &RolloutInfo{
							StartTime:       1643055431,
							EndTime:         1643055432,
							RolloutDuration: 1,
						},
						Output: &OutputVideo{
							PerceptualQuality: &OutputVideoPerceptualQualityVariableInline{
								Value: OutputVideoPerceptualQualityPtr(OutputVideoPerceptualQualityMediumHigh),
							},
						},
						Video:       true,
						User:        "jsmith",
						DateCreated: "2022-01-24 20:17:10+0000",
					},
				},
				TotalItems: 8,
			},
		},
		"200 OK very deep": {
			params: ListPoliciesRequest{
				Network:     PolicyNetworkStaging,
				Contract:    "3-WNKXX1",
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "itemKind": "POLICY",
    "items": [
{
	"variables": [
		{
			"name": "ResizeDim",
			"type": "number",
			"defaultValue": "280"
		},
		{
			"name": "ResizeDimWithBorder",
			"type": "number",
			"defaultValue": "260"
		},
		{
			"name": "MinDim",
			"type": "number",
			"defaultValue": "1000"
		},
		{
			"name": "MinDimNew",
			"type": "number",
			"defaultValue": "1450"
		},
		{
			"name": "MaxDimOld",
			"type": "number",
			"defaultValue": "1500"
		}
	],
	"transformations": [
		{
			"transformation": "Trim",
			"fuzz": 0.08,
			"padding": 0
		},
		{
			"transformation": "IfDimension",
			"dimension": "width",
			"value": {
				"var": "MaxDimOld"
			},
			"default": {
				"transformation": "Compound",
				"transformations": [
					{
						"transformation": "IfDimension",
						"dimension": "width",
						"value": {
							"var": "MinDim"
						},
						"lessThan": {
							"transformation": "Compound",
							"transformations": [
								{
									"transformation": "Resize",
									"aspect": "fit",
									"type": "normal",
									"width": {
										"var": "ResizeDimWithBorder"
									},
									"height": {
										"var": "ResizeDimWithBorder"
									}
								},
								{
									"transformation": "Crop",
									"xPosition": 0,
									"yPosition": 0,
									"gravity": "Center",
									"allowExpansion": true,
									"width": {
										"var": "ResizeDim"
									},
									"height": {
										"var": "ResizeDim"
									}
								},
								{
									"transformation": "BackgroundColor",
									"color": "#ffffff"
								}
							]
						},
						"default": {
							"transformation": "Compound",
							"transformations": [
								{
									"transformation": "IfDimension",
									"dimension": "height",
									"value": {
										"var": "MinDim"
									},
									"lessThan": {
										"transformation": "Compound",
										"transformations": [
											{
												"transformation": "Resize",
												"aspect": "fit",
												"type": "normal",
												"width": {
													"var": "ResizeDimWithBorder"
												},
												"height": {
													"var": "ResizeDimWithBorder"
												}
											},
											{
												"transformation": "Crop",
												"xPosition": 0,
												"yPosition": 0,
												"gravity": "Center",
												"allowExpansion": true,
												"width": {
													"var": "ResizeDim"
												},
												"height": {
													"var": "ResizeDim"
												}
											},
											{
												"transformation": "BackgroundColor",
												"color": "#ffffff"
											}
										]
									},
									"default": {
										"transformation": "Compound",
										"transformations": [
											{
												"transformation": "IfDimension",
												"dimension": "height",
												"value": {
													"var": "MaxDimOld"
												},
												"greaterThan": {
													"transformation": "Compound",
													"transformations": [
														{
															"transformation": "Resize",
															"aspect": "fit",
															"type": "normal",
															"width": {
																"var": "ResizeDimWithBorder"
															},
															"height": {
																"var": "ResizeDimWithBorder"
															}
														},
														{
															"transformation": "Crop",
															"xPosition": 0,
															"yPosition": 0,
															"gravity": "Center",
															"allowExpansion": true,
															"width": {
																"var": "ResizeDim"
															},
															"height": {
																"var": "ResizeDim"
															}
														},
														{
															"transformation": "BackgroundColor",
															"color": "#ffffff"
														}
													]
												},
												"default": {
													"transformation": "Compound",
													"transformations": [
														{
															"transformation": "Resize",
															"aspect": "fit",
															"type": "normal",
															"width": {
																"var": "ResizeDim"
															},
															"height": {
																"var": "ResizeDim"
															}
														},
														{
															"transformation": "Crop",
															"xPosition": 0,
															"yPosition": 0,
															"gravity": "Center",
															"allowExpansion": true,
															"width": {
																"var": "ResizeDim"
															},
															"height": {
																"var": "ResizeDim"
															}
														},
														{
															"transformation": "BackgroundColor",
															"color": "#ffffff"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				]
			}
		}
	],
	"breakpoints": {
		"widths": [
			280
		]
	},
	"output": {
		"perceptualQuality": "mediumHigh",
		"adaptiveQuality": 50
	},
	"video": false,
	"id": "multidimension",
	"dateCreated": "2022-01-01 12:00:00+0000",
	"previousVersion": 0,
	"version": 1
        }
    ],
    "totalItems": 8
}`,
			expectedPath: "/imaging/v2/network/staging/policies/",
			expectedResponse: &ListPoliciesResponse{
				ItemKind: "POLICY",
				Items: PolicyOutputs{
					&PolicyOutputImage{
						Variables: []Variable{
							{
								Name:         "ResizeDim",
								Type:         "number",
								DefaultValue: "280",
							},
							{
								Name:         "ResizeDimWithBorder",
								Type:         "number",
								DefaultValue: "260",
							},
							{
								Name:         "MinDim",
								Type:         "number",
								DefaultValue: "1000",
							},
							{
								Name:         "MinDimNew",
								Type:         "number",
								DefaultValue: "1450",
							},
							{
								Name:         "MaxDimOld",
								Type:         "number",
								DefaultValue: "1500",
							},
						},
						Transformations: []TransformationType{
							&Trim{
								Transformation: "Trim",
								Fuzz: &NumberVariableInline{
									Value: tools.Float32Ptr(0.08),
								},
								Padding: &IntegerVariableInline{
									Value: tools.IntPtr(0),
								},
							},
							&IfDimension{
								Transformation: "IfDimension",
								Dimension: &IfDimensionDimensionVariableInline{
									Value: IfDimensionDimensionPtr("width"),
								},
								Value: &IntegerVariableInline{
									Name: tools.StringPtr("MaxDimOld"),
								},
								Default: &Compound{
									Transformation: "Compound",
									Transformations: []TransformationType{
										&IfDimension{
											Transformation: "IfDimension",
											Dimension: &IfDimensionDimensionVariableInline{
												Value: IfDimensionDimensionPtr("width"),
											},
											Value: &IntegerVariableInline{
												Name: tools.StringPtr("MinDim"),
											},
											LessThan: &Compound{
												Transformation: "Compound",
												Transformations: []TransformationType{
													&Resize{
														Transformation: "Resize",
														Aspect: &ResizeAspectVariableInline{
															Value: ResizeAspectPtr("fit"),
														},
														Type: &ResizeTypeVariableInline{
															Value: ResizeTypePtr("normal"),
														},
														Width: &IntegerVariableInline{
															Name: tools.StringPtr("ResizeDimWithBorder"),
														},
														Height: &IntegerVariableInline{
															Name: tools.StringPtr("ResizeDimWithBorder"),
														},
													},
													&Crop{
														Transformation: "Crop",
														XPosition: &IntegerVariableInline{
															Value: tools.IntPtr(0),
														},
														YPosition: &IntegerVariableInline{
															Value: tools.IntPtr(0),
														},
														Gravity: &GravityVariableInline{
															Value: GravityPtr("Center"),
														},
														AllowExpansion: &BooleanVariableInline{
															Value: tools.BoolPtr(true),
														},
														Width: &IntegerVariableInline{
															Name: tools.StringPtr("ResizeDim"),
														},
														Height: &IntegerVariableInline{
															Name: tools.StringPtr("ResizeDim"),
														},
													},
													&BackgroundColor{
														Transformation: "BackgroundColor",
														Color: &StringVariableInline{
															Value: tools.StringPtr("#ffffff"),
														},
													},
												},
											},
											Default: &Compound{
												Transformation: "Compound",
												Transformations: []TransformationType{
													&IfDimension{
														Transformation: "IfDimension",
														Dimension: &IfDimensionDimensionVariableInline{
															Value: IfDimensionDimensionPtr("height"),
														},
														Value: &IntegerVariableInline{
															Name: tools.StringPtr("MinDim"),
														},
														LessThan: &Compound{
															Transformation: "Compound",
															Transformations: []TransformationType{
																&Resize{
																	Transformation: "Resize",
																	Aspect: &ResizeAspectVariableInline{
																		Value: ResizeAspectPtr("fit"),
																	},
																	Type: &ResizeTypeVariableInline{
																		Value: ResizeTypePtr("normal"),
																	},
																	Width: &IntegerVariableInline{
																		Name: tools.StringPtr("ResizeDimWithBorder"),
																	},
																	Height: &IntegerVariableInline{
																		Name: tools.StringPtr("ResizeDimWithBorder"),
																	},
																},
																&Crop{
																	Transformation: "Crop",
																	XPosition: &IntegerVariableInline{
																		Value: tools.IntPtr(0),
																	},
																	YPosition: &IntegerVariableInline{
																		Value: tools.IntPtr(0),
																	},
																	Gravity: &GravityVariableInline{
																		Value: GravityPtr("Center"),
																	},
																	AllowExpansion: &BooleanVariableInline{
																		Value: tools.BoolPtr(true),
																	},
																	Width: &IntegerVariableInline{
																		Name: tools.StringPtr("ResizeDim"),
																	},
																	Height: &IntegerVariableInline{
																		Name: tools.StringPtr("ResizeDim"),
																	},
																},
																&BackgroundColor{
																	Transformation: "BackgroundColor",
																	Color: &StringVariableInline{
																		Value: tools.StringPtr("#ffffff"),
																	},
																},
															},
														},
														Default: &Compound{
															Transformation: "Compound",
															Transformations: []TransformationType{
																&IfDimension{
																	Transformation: "IfDimension",
																	Dimension: &IfDimensionDimensionVariableInline{
																		Value: IfDimensionDimensionPtr("height"),
																	},
																	Value: &IntegerVariableInline{
																		Name: tools.StringPtr("MaxDimOld"),
																	},
																	GreaterThan: &Compound{
																		Transformation: "Compound",
																		Transformations: []TransformationType{
																			&Resize{
																				Transformation: "Resize",
																				Aspect: &ResizeAspectVariableInline{
																					Value: ResizeAspectPtr("fit"),
																				},
																				Type: &ResizeTypeVariableInline{
																					Value: ResizeTypePtr("normal"),
																				},

																				Width: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDimWithBorder"),
																				},
																				Height: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDimWithBorder"),
																				},
																			},
																			&Crop{
																				Transformation: "Crop",
																				XPosition: &IntegerVariableInline{
																					Value: tools.IntPtr(0),
																				},
																				YPosition: &IntegerVariableInline{
																					Value: tools.IntPtr(0),
																				},
																				Gravity: &GravityVariableInline{
																					Value: GravityPtr("Center"),
																				},
																				AllowExpansion: &BooleanVariableInline{
																					Value: tools.BoolPtr(true),
																				},
																				Width: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																				Height: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																			},
																			&BackgroundColor{
																				Transformation: "BackgroundColor",
																				Color: &StringVariableInline{
																					Value: tools.StringPtr("#ffffff"),
																				},
																			},
																		},
																	},
																	Default: &Compound{
																		Transformation: "Compound",
																		Transformations: []TransformationType{
																			&Resize{
																				Transformation: "Resize",
																				Aspect: &ResizeAspectVariableInline{
																					Value: ResizeAspectPtr("fit"),
																				},
																				Type: &ResizeTypeVariableInline{
																					Value: ResizeTypePtr("normal"),
																				},
																				Width: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																				Height: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																			},
																			&Crop{
																				Transformation: "Crop",
																				XPosition: &IntegerVariableInline{
																					Value: tools.IntPtr(0),
																				},
																				YPosition: &IntegerVariableInline{
																					Value: tools.IntPtr(0),
																				},
																				Gravity: &GravityVariableInline{
																					Value: GravityPtr("Center"),
																				},
																				AllowExpansion: &BooleanVariableInline{
																					Value: tools.BoolPtr(true),
																				},
																				Width: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																				Height: &IntegerVariableInline{
																					Name: tools.StringPtr("ResizeDim"),
																				},
																			},
																			&BackgroundColor{
																				Transformation: "BackgroundColor",
																				Color: &StringVariableInline{
																					Value: tools.StringPtr("#ffffff"),
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Breakpoints: &Breakpoints{
							Widths: []int{280},
						},
						Output: &OutputImage{
							PerceptualQuality: &OutputImagePerceptualQualityVariableInline{
								Value: OutputImagePerceptualQualityPtr("mediumHigh"),
							},
							AdaptiveQuality: 50,
						},
						Video:           false,
						ID:              "multidimension",
						DateCreated:     "2022-01-01 12:00:00+0000",
						PreviousVersion: 0,
						Version:         1,
					},
				},
				TotalItems: 8,
			},
		},
		"400 Bad request": {
			params: ListPoliciesRequest{
				Network:     PolicyNetworkStaging,
				Contract:    "3-WNKXX1",
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
"title": "Bad Request",
"instance": "52a21f40-9861-4d35-95d0-a603c85cb2ad",
"status": 400,
"detail": "A contract must be specified using the Contract header.",
"problemId": "52a21f40-9861-4d35-95d0-a603c85cb2ad"
}`,
			expectedPath: "/imaging/v2/network/staging/policies/",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
				Title:     "Bad Request",
				Instance:  "52a21f40-9861-4d35-95d0-a603c85cb2ad",
				Status:    400,
				Detail:    "A contract must be specified using the Contract header.",
				ProblemID: "52a21f40-9861-4d35-95d0-a603c85cb2ad",
			},
		},
		"401 Not authorized": {
			params: ListPoliciesRequest{
				Network:     PolicyNetworkStaging,
				Contract:    "3-WNKXX1",
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
"type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
"title": "Not authorized",
"status": 401,
"detail": "Inactive client token",
"instance": "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37c.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
"method": "GET",
"serverIp": "104.81.220.242",
"clientIp": "22.22.22.22",
"requestId": "124cc33c",
"requestTime": "2022-01-12T16:53:44Z"
}`,
			expectedPath: "/imaging/v2/network/staging/policies/",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37c.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "124cc33c",
				RequestTime: "2022-01-12T16:53:44Z",
			},
		},
		"403 Forbidden": {
			params: ListPoliciesRequest{
				Network:     PolicyNetworkStaging,
				Contract:    "3-WNKXX1",
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `{
				"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1002",
				"title": "Forbidden",
				"instance": "7d633d60-b120-4f28-a0de-ad86aeaf3c68",
				"status": 403,
				"detail": "User does not have authorization to perform this action.",
				"problemId": "7d633d60-b120-4f28-a0de-ad86aeaf3c68"
			}`,
			expectedPath: "/imaging/v2/network/staging/policies/",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1002",
				Title:     "Forbidden",
				Instance:  "7d633d60-b120-4f28-a0de-ad86aeaf3c68",
				Status:    403,
				Detail:    "User does not have authorization to perform this action.",
				ProblemID: "7d633d60-b120-4f28-a0de-ad86aeaf3c68",
			},
		},
		// 500
		"invalid network": {
			params: ListPoliciesRequest{
				Contract: "3-WNKXX1",
				Network:  "foo",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: ListPoliciesRequest{
				Network: PolicyNetworkProduction,
			},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				for h := range test.expectedHeaders {
					assert.Equal(t, test.expectedHeaders[h], r.Header[h])
				}
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListPolicies(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
