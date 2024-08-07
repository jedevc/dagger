defmodule Dagger.Codegen.ElixirGenerator.ObjectRendererTest do
  use ExUnit.Case, async: true
  use Mneme

  alias Dagger.Codegen.ElixirGenerator.ObjectRenderer

  test "return object node" do
    auto_assert(
      """
      # This file generated by `dagger_codegen`. Please DO NOT EDIT.
      defmodule Dagger.Client do
        @moduledoc "The root of the DAG."

        use Dagger.Core.QueryBuilder

        defstruct [:selection, :client]

        @type t() :: %__MODULE__{}

        @doc "Create a new TypeDef."
        @spec type_def(t()) :: Dagger.TypeDef.t()
        def type_def(%__MODULE__{} = client) do
          selection =
            client.selection |> select("typeDef")

          %Dagger.TypeDef{
            selection: selection,
            client: client.client
          }
        end
      end\
      """ <- render_type(ObjectRenderer, "test/fixtures/objects/chain-selection.json")
    )
  end

  test "return a list of leaf node" do
    auto_assert(
      """
      # This file generated by `dagger_codegen`. Please DO NOT EDIT.
      defmodule Dagger.Container do
        @moduledoc "The root of the DAG."

        use Dagger.Core.QueryBuilder

        defstruct [:selection, :client]

        @type t() :: %__MODULE__{}

        @doc "Retrieves the list of environment variables passed to commands."
        @spec env_variables(t()) :: {:ok, [Dagger.EnvVariable.t()]} | {:error, term()}
        def env_variables(%__MODULE__{} = container) do
          selection =
            container.selection |> select("envVariables") |> select("id")

          with {:ok, items} <- execute(selection, container.client) do
            {:ok,
             for %{"id" => id} <- items do
               %Dagger.EnvVariable{
                 selection:
                   query()
                   |> select("loadEnvVariableFromID")
                   |> arg("id", id),
                 client: container.client
               }
             end}
          end
        end
      end\
      """ <- render_type(ObjectRenderer, "test/fixtures/objects/list-leaf-nodes.json")
    )
  end

  test "execute leaf node" do
    auto_assert(
      """
      # This file generated by `dagger_codegen`. Please DO NOT EDIT.
      defmodule Dagger.EnvVariable do
        @moduledoc "An environment variable name and value."

        use Dagger.Core.QueryBuilder

        defstruct [:selection, :client]

        @type t() :: %__MODULE__{}

        @doc "The environment variable name."
        @spec name(t()) :: {:ok, String.t()} | {:error, term()}
        def name(%__MODULE__{} = env_variable) do
          selection =
            env_variable.selection |> select("name")

          execute(selection, env_variable.client)
        end
      end\
      """ <- render_type(ObjectRenderer, "test/fixtures/objects/execute-leaf-node.json")
    )
  end

  test "accept struct type of id argument" do
    auto_assert(
      """
      # This file generated by `dagger_codegen`. Please DO NOT EDIT.
      defmodule Dagger.Client do
        @moduledoc "The root of the DAG."

        use Dagger.Core.QueryBuilder

        defstruct [:selection, :client]

        @type t() :: %__MODULE__{}

        @doc "Load a Container from its ID."
        @spec load_container_from_id(t(), Dagger.ContainerID.t()) :: Dagger.Container.t()
        def load_container_from_id(%__MODULE__{} = client, id) do
          selection =
            client.selection |> select("loadContainerFromID") |> put_arg("id", id)

          %Dagger.Container{
            selection: selection,
            client: client.client
          }
        end
      end\
      """ <- render_type(ObjectRenderer, "test/fixtures/objects/id-arg.json")
    )
  end

  test "iss-7788 sync function return object instead of id type" do
    auto_assert(
      """
      # This file generated by `dagger_codegen`. Please DO NOT EDIT.
      defmodule Dagger.Container do
        @moduledoc "The root of the DAG."

        use Dagger.Core.QueryBuilder

        @derive Dagger.Sync
        defstruct [:selection, :client]

        @type t() :: %__MODULE__{}

        @doc \"""
        Forces evaluation of the pipeline in the engine.

        It doesn't run the default command if no exec has been set.
        \"""
        @spec sync(t()) :: {:ok, Dagger.Container.t()} | {:error, term()}
        def sync(%__MODULE__{} = container) do
          selection =
            container.selection |> select("sync")

          with {:ok, id} <- execute(selection, container.client) do
            {:ok,
             %Dagger.Container{
               selection:
                 query()
                 |> select("loadContainerFromID")
                 |> arg("id", id),
               client: container.client
             }}
          end
        end
      end\
      """ <- render_type(ObjectRenderer, "test/fixtures/objects/iss-7788.json")
    )
  end

  defp decode_type_from_file(path) do
    path
    |> File.read!()
    |> Jason.decode!()
    |> Dagger.Codegen.Introspection.Types.Type.from_map()
  end

  defp render(type, renderer) do
    renderer.render(type)
    |> IO.iodata_to_binary()
    |> Code.format_string!()
    |> IO.iodata_to_binary()
  end

  defp render_type(renderer, path) do
    path
    |> decode_type_from_file()
    |> render(renderer)
  end
end
