/*
 * This is a graph package, which contains graoh definition and basic operations
 * on it. As you go through the file, you will see some comments, that are
 * explaining this or that choice, etc.
 *
 * Author: github.com/tolstovrob
 */

package graph

type Option[T any] func(*T) // Type representing functional options pattern

type TKey uint64   // Key type. Can be replaced with any UNIQUE type
type TWeight int64 // Weight type. Can be replaced with any COMPARABLE type
