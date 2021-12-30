let data = require('./data.json')
let filtered = require("./filtered.json")
// let doubled = filtered.filter(e => {
//     // return e.ClassName.includes("MORETHAN1") || e.ClassName.includes("NOFOUND")
//     return e.exports.length > 2
// })
// let doubled = data.map(e => {
//     return e.exports.map(el => {
//         return el.data.properties.filter(ell => {
//             // l = ["ObjectProperty", "TextProperty", "EnumProperty", "ByteProperty", "NameProperty", "FloatProperty", "ArrayProperty", "IntProperty"]
//             // i = l.includes(ell.property_type.replace("\u0000", ""))
//             // return !i
//             return isArrayOfStructs(ell) || isStruct(ell)
//         })
//     })
// }).flat(99999)
debugger

const plist = [
    "ObjectProperty",
    "TextProperty",
    "EnumProperty",
    "ByteProperty",
    "NameProperty",
    "FloatProperty",
    "ArrayProperty",
    "IntProperty"
]


function isArrayOfStructs(e) {
    return isArrayOfType(e, "Struct")
}

function isArrayOfType(e, t) {
    return isArray(e)
        && isOfTypeStr(e.tag_data, t)
}

function isArray(e) {
    return isOfType(e, "Array")
}

function isStruct(e) {
    return isOfType(e, "Struct")
}

function isOfType(e, t) {
    return isOfTypeStr(e.property_type, t)
}

function isOfTypeStr(s, t) {
    return clean(s) === t + "Property"
}

function clean(s) {
    return s.replace("\u0000", "")
}