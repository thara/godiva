package class

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

func validate[T any](e *errReader, target T, vs ...validator[T]) {
	if e.err != nil {
		return
	}
	for _, v := range vs {
		if err := v.validate(target, e.name); err != nil {
			e.err = err
			return
		}
	}
}

type validator[T any] interface {
	validate(v T, name string) error
}

func match[T any](expected T) validator[T] {
	return &matchValidator[T]{expected: expected}
}

type matchValidator[T any] struct {
	expected T
}

func (v *matchValidator[T]) validate(target T, name string) error {
	if reflect.DeepEqual(v.expected, target) {
		return nil
	}
	return fmt.Errorf("%s does not match expected value (got:%v, want:%v)", name, target, v.expected)
}

func min[T constraints.Integer](minValue T) validator[T] {
	return &minValidator[T]{minValue: minValue}
}

type minValidator[T constraints.Integer] struct {
	minValue T
}

func (v *minValidator[T]) validate(target T, name string) error {
	if v.minValue <= target {
		return nil
	}
	return fmt.Errorf("%s out of range(got:%v, min:%d)", name, target, v.minValue)
}

func max[T constraints.Integer](maxValue T) validator[T] {
	return &maxValidator[T]{maxValue: maxValue}
}

type maxValidator[T constraints.Integer] struct {
	maxValue T
}

func (v *maxValidator[T]) validate(target T, name string) error {
	if target <= v.maxValue {
		return nil
	}
	return fmt.Errorf("%s out of range(got:%v, max:%d)", name, target, v.maxValue)
}

func existConstantPool[T constraints.Integer](cf *ClassFile) validator[T] {
	return &constantPoolExistanceValidator[T]{cp: cf.ConstantPool}
}

type constantPoolExistanceValidator[T constraints.Integer] struct {
	cp []cpInfo
}

func (v *constantPoolExistanceValidator[T]) validate(i T, name string) error {
	if 1 <= i && int(i) <= len(v.cp) {
		return nil
	}
	return fmt.Errorf("%s(%d) must be valid index in constant_pool", name, i)
}

func constantPoolStructure[T constraints.Integer, V cpInfo](cf *ClassFile) validator[T] {
	return &constantPoolStructureValidator[T, V]{cp: cf.ConstantPool}
}

type constantPoolStructureValidator[T constraints.Integer, V cpInfo] struct {
	cp []cpInfo
}

func (v *constantPoolStructureValidator[T, V]) validate(i T, name string) error {
	if i < 1 || len(v.cp) < int(i) {
		return fmt.Errorf("%s(%d) must be valid index in constant_pool", name, i)
	}

	entry := v.cp[i-1]
	_, ok := entry.(V)
	if ok {
		return nil
	}
	var e V
	return fmt.Errorf("constant_pool entry at `%s`(%d) must be a %s structure", name, i, e.StructureName())
}
