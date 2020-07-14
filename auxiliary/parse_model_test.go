package auxiliary

import (
	"fmt"
	"regexp"
	"testing"
)

func TestMatchModel_matchModule(t *testing.T) {
//	a := `XG 001/07AUG19      HGHMOW  ADT  RU03/GS /IPRFE  /327 /ATPCO/
//16.PENALTIES-CHANGES/CANCEL
//NOTE - RULE RU01 IN IPRG APPLIES
//<<         NOTE -
//<<          CANCELLATION
//<<          CHARGE CNY800 FOR REFUND/CANCEL
//<<          CHARGE CNY200 FOR NO-SHOW
//<<          CHILDREN DISCOUNT DO NOT APPLY
//<<          INFANT UNDER 2 WITHOUT A SEAT REFUNDS FREE
//<<         NOTE -
//<<          CHANGES
//<<          ANY TIME
//<<          CHARGE CNY500 FOR REISSUE/REVALIDATION
//<<          CHARGE CNY200 FOR NO-SHOW
//<<          ----------------------
//<<          CHILD CHARGES SAME AS ADULT.
//<<          INFANT FREE OF CHARGE
//<<          ----------------------
//<<          BEFORE DEPARTURE OF JOURNEY
//<<          REPRICE USING CURRENT FARES IN EFFECT AT THE
//<<          DATE OF TICKET REISSUANCE
//<<          ----------------------
//FSKY/1E/YGPH2Y1MVICY344/FCC=D/PAGE 1/2XG 001/07AUG19      HGHMOW  ADT  RU03/GS /IPRFE  /327 /ATPCO/
//<<          AFTER DEPARTURE OF JOURNEY
//<<          REPRICE USING HISTORICAL FARES IN EFFECT AT
//<<          ORIGINAL TICKET DATE.
//<<          ----------------------
//<<          CHARGE DIFFERENCE IN FARE PLUS CHANGE FEE IF ANY
//<<          ----------------------
//<<          IF THE NEW BASE FARE IS LOWER THAN THE PREVIOUS
//<<          BASE FARE IGNORE ANY RESIDUAL AMOUNT AND NO
//<<          REFUND PERMITTED.CHARGE CHANGE FEE IF ANY.
//<<          ----------------------
//<<          APPLY THE HIGHEST CHANGE FEE OF ALL CHANGED FARE
//<<          COMPONENTS WITHIN THE JOURNEY.
//<<          ----------------------
//<<          CHANGE FEE IS SUBJECT TO EACH CHANGE TRANSACTION
//<<          ----------------------
//<<          WAIVED FOR DEATH OF PASSENGER OR FAMILY MEMBER
//<<          DEATH CERTIFICATE REQUIRED.
//<<          ----------------------
//<<          CHANGES MUST BE WITHIN TICKET VALIDITY.
//FSKY/1E/YGPH2Y1MVICY344/FCC=D/PAGE 2/2`


//	b := ` CANCELLATION
//    <<          CHARGE CNY800 FOR REFUND/CANCEL
//    <<          CHARGE CNY200 FOR NO-SHOW`
	//reg := regexp.MustCompile(`CANCELLATION.*?<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR REFUND/CANCEL`)
	//reg := regexp.MustCompile(`CANCELLATION\s*<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR REFUND/CANCEL`)
	c := `XG 001/07AUG19      HGHMOW  ADT  RU03/GS /IPRFE  /327 /ATPCO/ 
16.PENALTIES-CHANGES/CANCEL 
 NOTE - RULE RU01 IN IPRG APPLIES   
<<         NOTE -   
<<          CANCELLATION
<<          CHARGE CNY800 FOR REFUND/CANCEL 
<<          CHARGE CNY200 FOR NO-SHOW   
<<          CHILDREN DISCOUNT DO NOT APPLY  
<<          INFANT UNDER 2 WITHOUT A SEAT REFUNDS FREE  
<<         NOTE -   
<<          CHANGES 
<<          ANY TIME
<<          CHARGE CNY500 FOR REISSUE/REVALIDATION  
<<          CHARGE CNY200 FOR NO-SHOW   
<<          ----------------------  
<<          CHILD CHARGES SAME AS ADULT.
<<          INFANT FREE OF CHARGE   
<<          ----------------------  
<<          BEFORE DEPARTURE OF JOURNEY 
<<          REPRICE USING CURRENT FARES IN EFFECT AT THE
<<          DATE OF TICKET REISSUANCE   
<<          ----------------------  
FSKY/1E/E3SQ3JCNVJK1A77/FCC=D/PAGE 1/2XG 001/07AUG19      HGHMOW  ADT  RU03/GS /IPRFE  /327 /ATPCO/ 
<<          AFTER DEPARTURE OF JOURNEY  
<<          REPRICE USING HISTORICAL FARES IN EFFECT AT 
<<          ORIGINAL TICKET DATE.   
<<          ----------------------  
<<          CHARGE DIFFERENCE IN FARE PLUS CHANGE FEE IF ANY
<<          ----------------------  
<<          IF THE NEW BASE FARE IS LOWER THAN THE PREVIOUS 
<<          BASE FARE IGNORE ANY RESIDUAL AMOUNT AND NO 
<<          REFUND PERMITTED.CHARGE CHANGE FEE IF ANY.  
<<          ----------------------  
<<          APPLY THE HIGHEST CHANGE FEE OF ALL CHANGED FARE
<<          COMPONENTS WITHIN THE JOURNEY.  
<<          ----------------------  
<<          CHANGE FEE IS SUBJECT TO EACH CHANGE TRANSACTION
<<          ----------------------  
<<          WAIVED FOR DEATH OF PASSENGER OR FAMILY MEMBER  
<<          DEATH CERTIFICATE REQUIRED. 
<<          ----------------------  
<<          CHANGES MUST BE WITHIN TICKET VALIDITY. 
FSKY/1E/E3SQ3JCNVJK1A77/FCC=D/PAGE 2/2`
	//c = ` <<          CANCELLATION
    //<<          CHARGE CNY800 FOR REFUND/CANCEL
    //<<          CHARGE CNY200 FOR NO-SHOW`
	//reg := regexp.MustCompile(`CANCELLATION\s<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR REFUND/CANCEL`)
	//reg := regexp.MustCompile(`CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR REFUND/CANCEL`)
	reg := regexp.MustCompile(`CANCELLATION\s<<.*?CHARGE.*?\s<<.*?CHARGE ([A-Z]{3})(\d+)(\.*\d*) FOR NO-SHOW`)
	subLst01 := reg.FindAllStringSubmatch(c, -1)
	fmt.Println(subLst01)


}