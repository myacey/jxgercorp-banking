function camel(str) {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase())
}

export function convertKeysToCamel(obj) {
  if (Array.isArray(obj)) {
    return obj.map(convertKeysToCamel)
  } else if (obj !== null && typeof obj === 'object') {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [
        camel(key),
        convertKeysToCamel(value)
      ])
    )
  }
  return obj
}
