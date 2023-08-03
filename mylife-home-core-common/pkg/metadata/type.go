package metadata

type Type interface {
}

/// <reference types="node" />
export interface Primitive {
	readonly id: string;
	encode(value: any): Buffer;
	decode(data: Buffer): any;
}
export declare namespace Primitives {
	const STRING: Primitive;
	const BOOL: Primitive;
	const UINT8: Primitive;
	const INT8: Primitive;
	const UINT32: Primitive;
	const INT32: Primitive;
	const FLOAT: Primitive;
	const JSON: Primitive;
}
export declare type TypeId = 'range' | 'text' | 'float' | 'bool' | 'enum' | 'complex';
export interface Type {
	readonly typeId: TypeId;
	readonly primitive: Primitive;
	toString(): string;
	validate(value: any): void;
}
export declare function parseType(value: string): Type;
export declare class Range implements Type {
	readonly min: number;
	readonly max: number;
	readonly primitive: Primitive;
	constructor(min: number, max: number);
	get typeId(): TypeId;
	toString(): string;
	validate(value: any): void;
}
export declare class Text implements Type {
	get typeId(): TypeId;
	get primitive(): Primitive;
	toString(): string;
	validate(value: any): void;
}
export declare class Float implements Type {
	get typeId(): TypeId;
	get primitive(): Primitive;
	toString(): string;
	validate(value: any): void;
}
export declare class Bool implements Type {
	get typeId(): TypeId;
	get primitive(): Primitive;
	toString(): string;
	validate(value: any): void;
}
export declare class Enum implements Type {
	readonly values: readonly string[];
	constructor(...values: string[]);
	get typeId(): TypeId;
	get primitive(): Primitive;
	toString(): string;
	validate(value: any): void;
}
export declare class Complex implements Type {
	get typeId(): TypeId;
	get primitive(): Primitive;
	toString(): string;
	validate(value: any): void;
}
