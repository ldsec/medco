# Survival analysis backend

## Inputs

Inputs are:

name | type | purpose
-- | -- | --
Start event concept path | string (with table access, formatted with slashes "/") | start event
End event concept path | string (with table access, formatted with slashes "/") | end event
Start event modifier code | string | start event
End event modifier code | string | end event
Sub group definitions | list of lists of i2b2 selection panels | group comparison
Relative time limit | integer (number of days) | greatest time point returned
Time granularity | string (one of day, week, month) | lower number of time points returned
Cohort ID | integer | analysis is run on that patient set
User ID | string | patient set retrieval
User public key | string (base64) | encryption of results

I2B2 selection panel is defined in [swagger definitions](https://github.com/ldsec/medco/connector/blob/merge-survival/swagger/medco-connector.yml#L571).

In changes to come, users will be able to indifferently use concept path or modifier path for events.

## Cohort

GetPatientList runs a SQL query that returns the list of patient numbers for given cohort ID and user name.

## Retrieve concept codes from paths

Survival query performs SQL queries directly on the observation fact by filtering concept and modifier codes. Concept codes are retrieved from their paths with the I2B2 ontology request GetTermInfo.

## Execute explore subqueries

Once concept codes are retrieved, one or multiple I2B2 PsmQuery and GetPatientSet requests are calledto find that have the start event concept.

### I2B2 requests

I2B2 PsmQuery returns an id that is passed as argument to I2B2 GetPatientSet to retreived the list of patient numbers.

### When no groups are provided

One I2B2 PsmQuery request is run. The parameters of this request is a list with one selection panel with one element. This element is the start event concept path.

### Panel redefinition

If sub groups are present, each list are appended one panel with the start event concept path.
For each sub group, a I2B2 PsmQuery request is called with the panel list of the sub group as argument.

## Cohort intersection

For any I2B2 results, the result patient number list is intersected with that of patient numbers returned by GetPatientList, so only patients who have start concept in observations and are in the cohort are kept.

## Build time points

The resulting patient lists are passed as input to SQL queries that build time points up to the given maximum limit. A time point is composed of three number:

* relative time
* number of events of interest \(defined by end event\)
* number of censoring events

### Relative times

Relative time is computed as the difference, in day, between the start date of end event and the start date of the end event.

### Handling right censoring events

Right censoring means that patients can leave the observations with an unkown status. Here, it is the case of a patient that has a the start event but not the end event. In this situation, the end event is the start date of the most recent observation.

## Granularity

The size of each group and all events, both censoring and of interest, are encrypted with ElGamal. To reduce the data transmission and the following decryption at client side, it is possible to change the resolution of relative points. The available granularities are day, week, month and year. Integers are ceiled: day 5 corresponds to week 1, month 1 and year 1.

## Expansion

For the moment, the following collective aggregation requires input arrays of same size across the nodes. This is ensured by adding all missing relative timees with  0 events of interest and 0 censoring events upto the max limit in day limit. With other granularities, the limit is the ceiling of \(limit in days divided by the number of day in the granularity\) +1.

## Encryption

Group size and number of events are encrypted with the collective authority key in elliptic curve ElGamal. These local aggregates are then passed to AKSgroups (Aggregate and Key Switch). AKSgroups, distributively aggregates the initial patient counts adn the events counts across nodes and switches the encryption key from the collective authority's to the user's.

## Outputs

The final results returned to client is a a set of time points group, one per sub group definition if these definitions are provided, otherwise one group. A group holds the encrypted initial number of patients, and the list of time points. One time point hsa three values, a relative time, a number of events of interest and a number of censoring events. The relative times are in plain text, the initial countsand  the two numbers of events are encrypted with the user's public key.
