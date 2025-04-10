{"Input":"(ab|aab|aaab)c", "Rv":1025, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"a", "Fr":0, "To":1 }]}
 { "Sn":1,  "Edge":[ { "On":"a", "Fr":1, "To":2 }, { "On":"b", "Fr":1, "To":3 }]}
 { "Sn":2,  "Edge":[ { "On":"a", "Fr":2, "To":4 }, { "On":"b", "Fr":2, "To":5 }]}
 { "Sn":3,  "Edge":[ { "On":"c", "Fr":3, "To":6 }]}
 { "Sn":4,  "Edge":[ { "On":"b", "Fr":4, "To":7 }]}
 { "Sn":5,  "Edge":[ { "On":"c", "Fr":5, "To":6 }]}
 { "Sn":6,  "Term":1025,  "Edge":[ ]}
 { "Sn":7,  "Edge":[ { "On":"c", "Fr":7, "To":6 }]}
]}
