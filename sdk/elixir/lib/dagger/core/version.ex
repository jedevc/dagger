defmodule Dagger.Core.Version do
  @moduledoc false

  @dagger_cli_version ""

  def engine_version(), do: @dagger_cli_version
end
