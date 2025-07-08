// This is a generated file - do not edit.
//
// Generated from x/customwasm/types/query.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class QueryParamsRequest extends $pb.GeneratedMessage {
  factory QueryParamsRequest() => create();

  QueryParamsRequest._();

  factory QueryParamsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.customwasm.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsRequest clone() => QueryParamsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsRequest copyWith(void Function(QueryParamsRequest) updates) => super.copyWith((message) => updates(message as QueryParamsRequest)) as QueryParamsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryParamsRequest create() => QueryParamsRequest._();
  @$core.override
  QueryParamsRequest createEmptyInstance() => create();
  static $pb.PbList<QueryParamsRequest> createRepeated() => $pb.PbList<QueryParamsRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryParamsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryParamsRequest>(create);
  static QueryParamsRequest? _defaultInstance;
}

class QueryParamsResponse extends $pb.GeneratedMessage {
  factory QueryParamsResponse({
    $fixnum.Int64? storeCodeBaseGas,
    $fixnum.Int64? storeCodeMultiplier,
    $fixnum.Int64? instBaseGas,
    $fixnum.Int64? instMultiplier,
  }) {
    final result = create();
    if (storeCodeBaseGas != null) result.storeCodeBaseGas = storeCodeBaseGas;
    if (storeCodeMultiplier != null) result.storeCodeMultiplier = storeCodeMultiplier;
    if (instBaseGas != null) result.instBaseGas = instBaseGas;
    if (instMultiplier != null) result.instMultiplier = instMultiplier;
    return result;
  }

  QueryParamsResponse._();

  factory QueryParamsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.customwasm.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'storeCodeBaseGas', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'storeCodeMultiplier', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'instBaseGas', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(4, _omitFieldNames ? '' : 'instMultiplier', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsResponse clone() => QueryParamsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsResponse copyWith(void Function(QueryParamsResponse) updates) => super.copyWith((message) => updates(message as QueryParamsResponse)) as QueryParamsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryParamsResponse create() => QueryParamsResponse._();
  @$core.override
  QueryParamsResponse createEmptyInstance() => create();
  static $pb.PbList<QueryParamsResponse> createRepeated() => $pb.PbList<QueryParamsResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryParamsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryParamsResponse>(create);
  static QueryParamsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get storeCodeBaseGas => $_getI64(0);
  @$pb.TagNumber(1)
  set storeCodeBaseGas($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasStoreCodeBaseGas() => $_has(0);
  @$pb.TagNumber(1)
  void clearStoreCodeBaseGas() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get storeCodeMultiplier => $_getI64(1);
  @$pb.TagNumber(2)
  set storeCodeMultiplier($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasStoreCodeMultiplier() => $_has(1);
  @$pb.TagNumber(2)
  void clearStoreCodeMultiplier() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get instBaseGas => $_getI64(2);
  @$pb.TagNumber(3)
  set instBaseGas($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasInstBaseGas() => $_has(2);
  @$pb.TagNumber(3)
  void clearInstBaseGas() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get instMultiplier => $_getI64(3);
  @$pb.TagNumber(4)
  set instMultiplier($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasInstMultiplier() => $_has(3);
  @$pb.TagNumber(4)
  void clearInstMultiplier() => $_clearField(4);
}

class QueryApi {
  final $pb.RpcClient _client;

  QueryApi(this._client);

  $async.Future<QueryParamsResponse> params($pb.ClientContext? ctx, QueryParamsRequest request) =>
    _client.invoke<QueryParamsResponse>(ctx, 'Query', 'Params', request, QueryParamsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
