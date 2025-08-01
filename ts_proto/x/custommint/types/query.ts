// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.7.5
//   protoc               v3.21.12
// source: x/custommint/types/query.proto

/* eslint-disable */
import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import Long from "long";
import { Coin } from "../../../cosmos/base/v1beta1/coin";

export const protobufPackage = "x.custommint.types";

export interface QueryTotalMintedRequest {
}

export interface QueryTotalMintedResponse {
  amount?: Coin | undefined;
}

export interface QueryTotalBurnedRequest {
}

export interface QueryTotalBurnedResponse {
  amount?: Coin | undefined;
}

export interface QueryParamsRequest {
}

export interface QueryParamsResponse {
  burnThreshold: string;
  burnRatio: string;
  blocksPerYear: Long;
  firstYearInflation: string;
  secondYearInflation: string;
  otherYearInflation: string;
}

function createBaseQueryTotalMintedRequest(): QueryTotalMintedRequest {
  return {};
}

export const QueryTotalMintedRequest: MessageFns<QueryTotalMintedRequest> = {
  encode(_: QueryTotalMintedRequest, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalMintedRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTotalMintedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): QueryTotalMintedRequest {
    return {};
  },

  toJSON(_: QueryTotalMintedRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryTotalMintedRequest>, I>>(base?: I): QueryTotalMintedRequest {
    return QueryTotalMintedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryTotalMintedRequest>, I>>(_: I): QueryTotalMintedRequest {
    const message = createBaseQueryTotalMintedRequest();
    return message;
  },
};

function createBaseQueryTotalMintedResponse(): QueryTotalMintedResponse {
  return { amount: undefined };
}

export const QueryTotalMintedResponse: MessageFns<QueryTotalMintedResponse> = {
  encode(message: QueryTotalMintedResponse, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalMintedResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTotalMintedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.amount = Coin.decode(reader, reader.uint32());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QueryTotalMintedResponse {
    return { amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined };
  },

  toJSON(message: QueryTotalMintedResponse): unknown {
    const obj: any = {};
    if (message.amount !== undefined) {
      obj.amount = Coin.toJSON(message.amount);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryTotalMintedResponse>, I>>(base?: I): QueryTotalMintedResponse {
    return QueryTotalMintedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryTotalMintedResponse>, I>>(object: I): QueryTotalMintedResponse {
    const message = createBaseQueryTotalMintedResponse();
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseQueryTotalBurnedRequest(): QueryTotalBurnedRequest {
  return {};
}

export const QueryTotalBurnedRequest: MessageFns<QueryTotalBurnedRequest> = {
  encode(_: QueryTotalBurnedRequest, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalBurnedRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTotalBurnedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): QueryTotalBurnedRequest {
    return {};
  },

  toJSON(_: QueryTotalBurnedRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryTotalBurnedRequest>, I>>(base?: I): QueryTotalBurnedRequest {
    return QueryTotalBurnedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryTotalBurnedRequest>, I>>(_: I): QueryTotalBurnedRequest {
    const message = createBaseQueryTotalBurnedRequest();
    return message;
  },
};

function createBaseQueryTotalBurnedResponse(): QueryTotalBurnedResponse {
  return { amount: undefined };
}

export const QueryTotalBurnedResponse: MessageFns<QueryTotalBurnedResponse> = {
  encode(message: QueryTotalBurnedResponse, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.amount !== undefined) {
      Coin.encode(message.amount, writer.uint32(10).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryTotalBurnedResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTotalBurnedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.amount = Coin.decode(reader, reader.uint32());
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QueryTotalBurnedResponse {
    return { amount: isSet(object.amount) ? Coin.fromJSON(object.amount) : undefined };
  },

  toJSON(message: QueryTotalBurnedResponse): unknown {
    const obj: any = {};
    if (message.amount !== undefined) {
      obj.amount = Coin.toJSON(message.amount);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryTotalBurnedResponse>, I>>(base?: I): QueryTotalBurnedResponse {
    return QueryTotalBurnedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryTotalBurnedResponse>, I>>(object: I): QueryTotalBurnedResponse {
    const message = createBaseQueryTotalBurnedResponse();
    message.amount = (object.amount !== undefined && object.amount !== null)
      ? Coin.fromPartial(object.amount)
      : undefined;
    return message;
  },
};

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest: MessageFns<QueryParamsRequest> = {
  encode(_: QueryParamsRequest, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(base?: I): QueryParamsRequest {
    return QueryParamsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    burnThreshold: "",
    burnRatio: "",
    blocksPerYear: Long.ZERO,
    firstYearInflation: "",
    secondYearInflation: "",
    otherYearInflation: "",
  };
}

export const QueryParamsResponse: MessageFns<QueryParamsResponse> = {
  encode(message: QueryParamsResponse, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.burnThreshold !== "") {
      writer.uint32(10).string(message.burnThreshold);
    }
    if (message.burnRatio !== "") {
      writer.uint32(18).string(message.burnRatio);
    }
    if (!message.blocksPerYear.equals(Long.ZERO)) {
      writer.uint32(24).int64(message.blocksPerYear.toString());
    }
    if (message.firstYearInflation !== "") {
      writer.uint32(34).string(message.firstYearInflation);
    }
    if (message.secondYearInflation !== "") {
      writer.uint32(42).string(message.secondYearInflation);
    }
    if (message.otherYearInflation !== "") {
      writer.uint32(50).string(message.otherYearInflation);
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    const end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.burnThreshold = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.burnRatio = reader.string();
          continue;
        }
        case 3: {
          if (tag !== 24) {
            break;
          }

          message.blocksPerYear = Long.fromString(reader.int64().toString());
          continue;
        }
        case 4: {
          if (tag !== 34) {
            break;
          }

          message.firstYearInflation = reader.string();
          continue;
        }
        case 5: {
          if (tag !== 42) {
            break;
          }

          message.secondYearInflation = reader.string();
          continue;
        }
        case 6: {
          if (tag !== 50) {
            break;
          }

          message.otherYearInflation = reader.string();
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return {
      burnThreshold: isSet(object.burnThreshold) ? globalThis.String(object.burnThreshold) : "",
      burnRatio: isSet(object.burnRatio) ? globalThis.String(object.burnRatio) : "",
      blocksPerYear: isSet(object.blocksPerYear) ? Long.fromValue(object.blocksPerYear) : Long.ZERO,
      firstYearInflation: isSet(object.firstYearInflation) ? globalThis.String(object.firstYearInflation) : "",
      secondYearInflation: isSet(object.secondYearInflation) ? globalThis.String(object.secondYearInflation) : "",
      otherYearInflation: isSet(object.otherYearInflation) ? globalThis.String(object.otherYearInflation) : "",
    };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    if (message.burnThreshold !== "") {
      obj.burnThreshold = message.burnThreshold;
    }
    if (message.burnRatio !== "") {
      obj.burnRatio = message.burnRatio;
    }
    if (!message.blocksPerYear.equals(Long.ZERO)) {
      obj.blocksPerYear = (message.blocksPerYear || Long.ZERO).toString();
    }
    if (message.firstYearInflation !== "") {
      obj.firstYearInflation = message.firstYearInflation;
    }
    if (message.secondYearInflation !== "") {
      obj.secondYearInflation = message.secondYearInflation;
    }
    if (message.otherYearInflation !== "") {
      obj.otherYearInflation = message.otherYearInflation;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(base?: I): QueryParamsResponse {
    return QueryParamsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.burnThreshold = object.burnThreshold ?? "";
    message.burnRatio = object.burnRatio ?? "";
    message.blocksPerYear = (object.blocksPerYear !== undefined && object.blocksPerYear !== null)
      ? Long.fromValue(object.blocksPerYear)
      : Long.ZERO;
    message.firstYearInflation = object.firstYearInflation ?? "";
    message.secondYearInflation = object.secondYearInflation ?? "";
    message.otherYearInflation = object.otherYearInflation ?? "";
    return message;
  },
};

export interface Query {
  TotalMinted(request: QueryTotalMintedRequest): Promise<QueryTotalMintedResponse>;
  TotalBurned(request: QueryTotalBurnedRequest): Promise<QueryTotalBurnedResponse>;
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
}

export const QueryServiceName = "x.custommint.types.Query";
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || QueryServiceName;
    this.rpc = rpc;
    this.TotalMinted = this.TotalMinted.bind(this);
    this.TotalBurned = this.TotalBurned.bind(this);
    this.Params = this.Params.bind(this);
  }
  TotalMinted(request: QueryTotalMintedRequest): Promise<QueryTotalMintedResponse> {
    const data = QueryTotalMintedRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "TotalMinted", data);
    return promise.then((data) => QueryTotalMintedResponse.decode(new BinaryReader(data)));
  }

  TotalBurned(request: QueryTotalBurnedRequest): Promise<QueryTotalBurnedResponse> {
    const data = QueryTotalBurnedRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "TotalBurned", data);
    return promise.then((data) => QueryTotalBurnedResponse.decode(new BinaryReader(data)));
  }

  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new BinaryReader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Long ? string | number | Long : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export interface MessageFns<T> {
  encode(message: T, writer?: BinaryWriter): BinaryWriter;
  decode(input: BinaryReader | Uint8Array, length?: number): T;
  fromJSON(object: any): T;
  toJSON(message: T): unknown;
  create<I extends Exact<DeepPartial<T>, I>>(base?: I): T;
  fromPartial<I extends Exact<DeepPartial<T>, I>>(object: I): T;
}
