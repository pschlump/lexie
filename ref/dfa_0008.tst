{"Input":"(aa|bb|ccc)*abb", "Rv":1008, "Start": 0, "States":[
 { "Sn":0,  "Edge":[ { "On":"a", "Fr":0, "To":1 }, { "On":"b", "Fr":0, "To":2 }, { "On":"c", "Fr":0, "To":3 }]}
 { "Sn":1,  "Edge":[ { "On":"a", "Fr":1, "To":4 }, { "On":"b", "Fr":1, "To":5 }]}
 { "Sn":2,  "Edge":[ { "On":"b", "Fr":2, "To":6 }]}
 { "Sn":3,  "Edge":[ { "On":"c", "Fr":3, "To":7 }]}
 { "Sn":4,  "Edge":[ { "On":"a", "Fr":4, "To":1 }, { "On":"b", "Fr":4, "To":2 }, { "On":"c", "Fr":4, "To":3 }]}
 { "Sn":5,  "Edge":[ { "On":"b", "Fr":5, "To":8 }]}
 { "Sn":6,  "Edge":[ { "On":"a", "Fr":6, "To":1 }, { "On":"b", "Fr":6, "To":2 }, { "On":"c", "Fr":6, "To":3 }]}
 { "Sn":7,  "Edge":[ { "On":"c", "Fr":7, "To":9 }]}
 { "Sn":8,  "Term":1008,  "Edge":[ ]}
 { "Sn":9,  "Edge":[ { "On":"a", "Fr":9, "To":1 }, { "On":"b", "Fr":9, "To":2 }, { "On":"c", "Fr":9, "To":3 }]}
]}
