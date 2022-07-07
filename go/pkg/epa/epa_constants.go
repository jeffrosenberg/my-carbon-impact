package epa

// This file contains EPA constants for calculating carbon emissions
// Source: https://www3.epa.gov/carbon-footprint-calculator/

//
// Vehicle data
//
const VEHICLE_LBS_PER_GAL = 19.6           // 2011
const VEHICLE_OTHER_EMISSIONS_RATIO = 1.01 // 2011
const VEHICLE_AVG_MILES_PER_YEAR = 11398   // 2012
const VEHICLE_AVG_FUEL_ECONOMY_MPG = 21.6  // 2012
var VEHICLE_AVG_EMISSIONS = (VEHICLE_AVG_MILES_PER_YEAR / VEHICLE_AVG_FUEL_ECONOMY_MPG) * VEHICLE_LBS_PER_GAL * VEHICLE_OTHER_EMISSIONS_RATIO
