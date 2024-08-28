# This file generated by `dagger_codegen`. Please DO NOT EDIT.
defmodule Dagger.ModuleSource do
  @moduledoc "The source needed to load and run a module, along with any metadata about the source such as versions/urls/etc."

  alias Dagger.Core.Client
  alias Dagger.Core.QueryBuilder, as: QB

  @derive Dagger.ID

  defstruct [:query_builder, :client]

  @type t() :: %__MODULE__{}

  @doc "If the source is a of kind git, the git source representation of it."
  @spec as_git_source(t()) :: Dagger.GitModuleSource.t() | nil
  def as_git_source(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("asGitSource")

    %Dagger.GitModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "If the source is of kind local, the local source representation of it."
  @spec as_local_source(t()) :: Dagger.LocalModuleSource.t() | nil
  def as_local_source(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("asLocalSource")

    %Dagger.LocalModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Load the source as a module. If this is a local source, the parent directory must have been provided during module source creation"
  @spec as_module(t(), [{:engine_version, String.t() | nil}]) :: Dagger.Module.t()
  def as_module(%__MODULE__{} = module_source, optional_args \\ []) do
    query_builder =
      module_source.query_builder
      |> QB.select("asModule")
      |> QB.maybe_put_arg("engineVersion", optional_args[:engine_version])

    %Dagger.Module{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "A human readable ref string representation of this module source."
  @spec as_string(t()) :: {:ok, String.t()} | {:error, term()}
  def as_string(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("asString")

    Client.execute(module_source.client, query_builder)
  end

  @doc "Returns whether the module source has a configuration file."
  @spec config_exists(t()) :: {:ok, boolean()} | {:error, term()}
  def config_exists(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("configExists")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The directory containing everything needed to load load and use the module."
  @spec context_directory(t()) :: Dagger.Directory.t()
  def context_directory(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("contextDirectory")

    %Dagger.Directory{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "The dependencies of the module source. Includes dependencies from the configuration and any extras from withDependencies calls."
  @spec dependencies(t()) :: {:ok, [Dagger.ModuleDependency.t()]} | {:error, term()}
  def dependencies(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("dependencies") |> QB.select("id")

    with {:ok, items} <- Client.execute(module_source.client, query_builder) do
      {:ok,
       for %{"id" => id} <- items do
         %Dagger.ModuleDependency{
           query_builder:
             QB.query()
             |> QB.select("loadModuleDependencyFromID")
             |> QB.put_arg("id", id),
           client: module_source.client
         }
       end}
    end
  end

  @doc "Return the module source's content digest. The format of the digest is not guaranteed to be stable between releases of Dagger. It is guaranteed to be stable between invocations of the same Dagger engine."
  @spec digest(t()) :: {:ok, String.t()} | {:error, term()}
  def digest(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("digest")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The directory containing the module configuration and source code (source code may be in a subdir)."
  @spec directory(t(), String.t()) :: Dagger.Directory.t()
  def directory(%__MODULE__{} = module_source, path) do
    query_builder =
      module_source.query_builder |> QB.select("directory") |> QB.put_arg("path", path)

    %Dagger.Directory{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "A unique identifier for this ModuleSource."
  @spec id(t()) :: {:ok, Dagger.ModuleSourceID.t()} | {:error, term()}
  def id(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("id")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The kind of source (e.g. local, git, etc.)"
  @spec kind(t()) :: {:ok, Dagger.ModuleSourceKind.t()} | {:error, term()}
  def kind(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("kind")

    case Client.execute(module_source.client, query_builder) do
      {:ok, enum} -> {:ok, Dagger.ModuleSourceKind.from_string(enum)}
      error -> error
    end
  end

  @doc "If set, the name of the module this source references, including any overrides at runtime by callers."
  @spec module_name(t()) :: {:ok, String.t()} | {:error, term()}
  def module_name(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("moduleName")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The original name of the module this source references, as defined in the module configuration."
  @spec module_original_name(t()) :: {:ok, String.t()} | {:error, term()}
  def module_original_name(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("moduleOriginalName")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The path to the module source's context directory on the caller's filesystem. Only valid for local sources."
  @spec resolve_context_path_from_caller(t()) :: {:ok, String.t()} | {:error, term()}
  def resolve_context_path_from_caller(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("resolveContextPathFromCaller")

    Client.execute(module_source.client, query_builder)
  end

  @doc "Resolve the provided module source arg as a dependency relative to this module source."
  @spec resolve_dependency(t(), Dagger.ModuleSource.t()) :: Dagger.ModuleSource.t()
  def resolve_dependency(%__MODULE__{} = module_source, dep) do
    query_builder =
      module_source.query_builder
      |> QB.select("resolveDependency")
      |> QB.put_arg("dep", Dagger.ID.id!(dep))

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Load a directory from the caller optionally with a given view applied."
  @spec resolve_directory_from_caller(t(), String.t(), [{:view_name, String.t() | nil}]) ::
          Dagger.Directory.t()
  def resolve_directory_from_caller(%__MODULE__{} = module_source, path, optional_args \\ []) do
    query_builder =
      module_source.query_builder
      |> QB.select("resolveDirectoryFromCaller")
      |> QB.put_arg("path", path)
      |> QB.maybe_put_arg("viewName", optional_args[:view_name])

    %Dagger.Directory{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Load the source from its path on the caller's filesystem, including only needed+configured files and directories. Only valid for local sources."
  @spec resolve_from_caller(t()) :: Dagger.ModuleSource.t()
  def resolve_from_caller(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("resolveFromCaller")

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "The path relative to context of the root of the module source, which contains dagger.json. It also contains the module implementation source code, but that may or may not being a subdir of this root."
  @spec source_root_subpath(t()) :: {:ok, String.t()} | {:error, term()}
  def source_root_subpath(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("sourceRootSubpath")

    Client.execute(module_source.client, query_builder)
  end

  @doc "The path relative to context of the module implementation source code."
  @spec source_subpath(t()) :: {:ok, String.t()} | {:error, term()}
  def source_subpath(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("sourceSubpath")

    Client.execute(module_source.client, query_builder)
  end

  @doc "Retrieve a named view defined for this module source."
  @spec view(t(), String.t()) :: Dagger.ModuleSourceView.t()
  def view(%__MODULE__{} = module_source, name) do
    query_builder =
      module_source.query_builder |> QB.select("view") |> QB.put_arg("name", name)

    %Dagger.ModuleSourceView{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "The named views defined for this module source, which are sets of directory filters that can be applied to directory arguments provided to functions."
  @spec views(t()) :: {:ok, [Dagger.ModuleSourceView.t()]} | {:error, term()}
  def views(%__MODULE__{} = module_source) do
    query_builder =
      module_source.query_builder |> QB.select("views") |> QB.select("id")

    with {:ok, items} <- Client.execute(module_source.client, query_builder) do
      {:ok,
       for %{"id" => id} <- items do
         %Dagger.ModuleSourceView{
           query_builder:
             QB.query()
             |> QB.select("loadModuleSourceViewFromID")
             |> QB.put_arg("id", id),
           client: module_source.client
         }
       end}
    end
  end

  @doc "Update the module source with a new context directory. Only valid for local sources."
  @spec with_context_directory(t(), Dagger.Directory.t()) :: Dagger.ModuleSource.t()
  def with_context_directory(%__MODULE__{} = module_source, dir) do
    query_builder =
      module_source.query_builder
      |> QB.select("withContextDirectory")
      |> QB.put_arg("dir", Dagger.ID.id!(dir))

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Append the provided dependencies to the module source's dependency list."
  @spec with_dependencies(t(), [Dagger.ModuleDependencyID.t()]) :: Dagger.ModuleSource.t()
  def with_dependencies(%__MODULE__{} = module_source, dependencies) do
    query_builder =
      module_source.query_builder
      |> QB.select("withDependencies")
      |> QB.put_arg("dependencies", dependencies)

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Sets module init arguments"
  @spec with_init(t(), [{:merge, boolean() | nil}]) :: Dagger.ModuleSource.t()
  def with_init(%__MODULE__{} = module_source, optional_args \\ []) do
    query_builder =
      module_source.query_builder
      |> QB.select("withInit")
      |> QB.maybe_put_arg("merge", optional_args[:merge])

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Update the module source with a new name."
  @spec with_name(t(), String.t()) :: Dagger.ModuleSource.t()
  def with_name(%__MODULE__{} = module_source, name) do
    query_builder =
      module_source.query_builder |> QB.select("withName") |> QB.put_arg("name", name)

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Update the module source with a new SDK."
  @spec with_sdk(t(), String.t()) :: Dagger.ModuleSource.t()
  def with_sdk(%__MODULE__{} = module_source, sdk) do
    query_builder =
      module_source.query_builder |> QB.select("withSDK") |> QB.put_arg("sdk", sdk)

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Update the module source with a new source subpath."
  @spec with_source_subpath(t(), String.t()) :: Dagger.ModuleSource.t()
  def with_source_subpath(%__MODULE__{} = module_source, path) do
    query_builder =
      module_source.query_builder |> QB.select("withSourceSubpath") |> QB.put_arg("path", path)

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end

  @doc "Update the module source with a new named view."
  @spec with_view(t(), String.t(), [String.t()]) :: Dagger.ModuleSource.t()
  def with_view(%__MODULE__{} = module_source, name, patterns) do
    query_builder =
      module_source.query_builder
      |> QB.select("withView")
      |> QB.put_arg("name", name)
      |> QB.put_arg("patterns", patterns)

    %Dagger.ModuleSource{
      query_builder: query_builder,
      client: module_source.client
    }
  end
end
